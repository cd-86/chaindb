package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shynur/chaindb/blockchain"
	"github.com/shynur/chaindb/chaindb_config"
	"github.com/shynur/chaindb/discovery"
)

func StartUserService() {
	server := http.NewServeMux()

	var wait_for_sync sync.RWMutex
	go func() {
		for ; ; time.Sleep(chaindb_config.DiscoveryInterval) {
			func() {
				wait_for_sync.Lock()
				defer wait_for_sync.Unlock()

				peers := *discovery.ActiveNodes.Load()
				if len(peers) == 0 {
					return
				}

				h_heads := make(chan uint32, len(peers))
				for _, peer := range peers {
					go func() {
						head, err := func() (uint64, error) {
							resp, err := http.Get(
								fmt.Sprintf(
									"http://%s/head",
									net.JoinHostPort(
										peer,
										strconv.Itoa(chaindb_config.TCPPortMonitor),
									),
								),
							)
							if err != nil {
								return 0, err
							} else {
								defer resp.Body.Close()
							}

							body, err := io.ReadAll(resp.Body)
							if err != nil {
								return 0, err
							}

							var head uint64
							err = json.Unmarshal(body, &head)
							if err != nil {
								return 0, err
							}
							return head, nil
						}()
						if err != nil {
							return
						}

						resp, err := http.Get(
							fmt.Sprintf(
								"http://%s/blocks/%d",
								net.JoinHostPort(
									peer,
									strconv.Itoa(chaindb_config.TCPPortMonitor),
								),
								head,
							),
						)
						if err != nil {
							return
						} else {
							defer resp.Body.Close()
						}

						var blk blockchain.Block
						err = json.NewDecoder(resp.Body).Decode(&blk)
						if err != nil {
							return
						}
						h_heads <- blk.Height
					}()
				}

				highest := LocalChain.Head().Height
				timeout := time.After(chaindb_config.BlockTime)
				for c := true; c; {
					select {
					case <-timeout:
						c = false
					case h := <-h_heads:
						highest = max(highest, h)
					}
				}
				for highest > LocalChain.Head().Height+chaindb_config.MaxAllowedHeightReduction {
					time.Sleep(2 * chaindb_config.BlockTime)
				}
			}()
		}
	}()

	get_network_ip := func(connectible string) (string, error) {
		addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(connectible, "53"))
		if err != nil {
			return "", err
		}
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			return "", err
		}
		defer conn.Close()
		return conn.LocalAddr().(*net.UDPAddr).IP.String(), nil
	}

	// `DELETE /api/v1/peers/[host]`
	// ==> 状态码
	server.HandleFunc("/api/v1/peers/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "只支持 DELETE", http.StatusMethodNotAllowed)
			return
		}

		peer := strings.Split(r.URL.Path, "/")[4]

		this_ip, err := get_network_ip(peer)
		if err != nil {
			http.Error(w, "获取服务器自身 IP 失败", http.StatusInternalServerError)
			return
		}

		_, err = http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf(
				"http://%s/peers/%s",
				net.JoinHostPort(
					peer,
					strconv.Itoa(chaindb_config.TCPPortMonitor),
				),
				this_ip,
			),
			nil,
		)
		if err != nil {
			log.Printf(
				"[  User ] [Peer](%s) 删除 [本节点](%s): %v\n",
				peer, this_ip,
				err,
			)
		} else {
			log.Printf(
				"[  User ] 让 [Peer](%s) 删除 [本节点](%s)\n",
				peer, this_ip,
			)
		}

		http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf(
				"http://localhost:%d/peers/%s",
				chaindb_config.TCPPortMonitor,
				peer,
			),
			nil,
		)

		w.WriteHeader(http.StatusNoContent)
	})

	server.HandleFunc("/api/v1/peers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// `GET /api/v1/peers`
			// ==> JSON

			resp, err := http.Get(
				fmt.Sprintf(
					"http://localhost:%d/peers",
					chaindb_config.TCPPortMonitor,
				),
			)
			if err != nil {
				http.Error(w, "获取 peers 列表失败", http.StatusBadGateway)
				return
			} else {
				defer resp.Body.Close()
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
		case http.MethodPost:
			// `POST /api/v1/peers` + JSON
			// ==> 状态码

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "读取 HTTP Body 失败", http.StatusBadRequest)
				return
			}
			var peer string
			if err := json.Unmarshal(body, &peer); err != nil {
				http.Error(w, "Peer JSON 反序列化失败", http.StatusBadRequest)
				return
			}

			this_ip, err := get_network_ip(peer)
			if err != nil {
				http.Error(w, "获取服务器自身 IP 失败", http.StatusInternalServerError)
				return
			}
			_, err = http.Post(
				fmt.Sprintf(
					"http://%s/peers",
					net.JoinHostPort(
						peer,
						strconv.Itoa(chaindb_config.TCPPortMonitor),
					),
				),
				"application/json",
				bytes.NewReader(func() []byte {
					ip_json, _ := json.Marshal([]string{this_ip})
					return ip_json
				}()),
			)
			if err != nil {
				http.Error(w, "添加 Peer 失败", http.StatusBadGateway)
				return
			} else {
				log.Printf(
					"[  User ] 向 [Peer](%s) 添加 [本节点](%s)\n",
					peer, this_ip,
				)
			}

			http.Post(
				fmt.Sprintf(
					"http://localhost:%d/peers",
					chaindb_config.TCPPortMonitor,
				),
				"application/json",
				bytes.NewReader(func() []byte {
					ip_json, _ := json.Marshal([]string{peer})
					return ip_json
				}()),
			)

			w.WriteHeader(http.StatusCreated)
		default:
			http.Error(w, "不支持的 HTTP verb", http.StatusMethodNotAllowed)
		}
	})

	server.HandleFunc("/api/v1/transactions", func(w http.ResponseWriter, r *http.Request) {
		wait_for_sync.RLock()
		defer wait_for_sync.RUnlock()

		switch r.Method {
		case http.MethodPost:
			// `POST /api/v1/transactions?confirmation=6` + JSON
			// ==> 状态码

			if r.URL.Query().Get("confirmation") == "" {
				http.Error(w, "缺少 'confirmation' 参数", http.StatusBadRequest)
				return
			}
			confirmation, _ := strconv.Atoi(r.URL.Query().Get("confirmation"))

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "读取 HTTP Body 失败", http.StatusBadRequest)
				return
			}
			var tx_to_insert blockchain.Transaction
			if err := json.Unmarshal(body, &tx_to_insert); err != nil {
				http.Error(
					w,
					fmt.Sprintf("Tx JSON (%+v) 反序列化失败: %v", string(body), err),
					http.StatusBadRequest,
				)
				return
			}

			stop := make(chan any, 1)
			go func() {
				for ; ; time.Sleep(chaindb_config.BlockTime / 2) {
					select {
					default:
					case <-stop:
						return
					}

					active_peers := *discovery.ActiveNodes.Load()
					if len(active_peers) == 0 {
						http.Post(
							fmt.Sprintf(
								"http://localhost:%d/transactions",
								chaindb_config.TCPPortMonitor,
							),
							"application/json",
							bytes.NewReader(body),
						)
					} else {
						for _, peer := range active_peers {
							http.Post(
								fmt.Sprintf(
									"http://%s/transactions",
									net.JoinHostPort(
										peer,
										strconv.Itoa(chaindb_config.TCPPortMonitor),
									),
								),
								"application/json",
								bytes.NewReader(body),
							)
						}
					}
				}
			}()

			select {
			case inserted_tx := <-func() <-chan blockchain.Transaction {
				confirmation_fullfilled_tx := make(chan blockchain.Transaction, 1)

				go func() {
					for ; ; time.Sleep(chaindb_config.BlockTime) {
						resp, err := http.Get(
							fmt.Sprintf(
								"http://localhost:%d/owners/%d/transactions",
								chaindb_config.TCPPortMonitor,
								tx_to_insert.OwnerID,
							),
						)
						if err != nil {
							continue
						}

						var txs []struct {
							Tx                blockchain.Transaction
							ConfirmationScore uint32
						}
						err = json.NewDecoder(resp.Body).Decode(&txs)
						resp.Body.Close()
						if err != nil {
							continue
						}

						for _, tx := range txs {
							if tx.Tx.Nonce == tx_to_insert.Nonce &&
								tx.ConfirmationScore >= uint32(confirmation) {
								confirmation_fullfilled_tx <- tx.Tx
								return
							}
						}
					}
				}()

				return confirmation_fullfilled_tx
			}():
				if tx_to_insert.Data == inserted_tx.Data {
					w.WriteHeader(http.StatusCreated)
				} else {
					http.Error(w, "Tx 已经存在", http.StatusConflict)
				}
			case <-time.After(time.Duration(3 * (confirmation + 1) * int(chaindb_config.BlockTime))):
				http.Error(w, "超时", http.StatusServiceUnavailable)
			}
			stop <- nil
		case http.MethodGet:
			// `GET /api/v1/transactions?owner=42`
			// ==> JSON

			owner := r.URL.Query().Get("owner")
			if owner == "" {
				http.Error(w, "未指定 'owner'", http.StatusBadRequest)
				return
			}

			resp, err := http.Get(
				fmt.Sprintf(
					"http://localhost:%d/owners/%s/transactions",
					chaindb_config.TCPPortMonitor,
					owner,
				),
			)
			if err != nil {
				http.Error(w, "获取 owner 的数据列表失败", http.StatusBadGateway)
				return
			} else {
				defer resp.Body.Close()
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
		default:
			http.Error(w, "不支持的 HTTP verb", http.StatusMethodNotAllowed)
		}
	})

	// `GET /api/v1/owners`
	// ==> JSON
	server.HandleFunc("/api/v1/owners", func(w http.ResponseWriter, r *http.Request) {
		wait_for_sync.RLock()
		defer wait_for_sync.RUnlock()

		if r.Method != http.MethodGet {
			http.Error(w, "只支持 GET", http.StatusMethodNotAllowed)
			return
		}

		resp, err := http.Get(
			fmt.Sprintf(
				"http://localhost:%d/owners",
				chaindb_config.TCPPortMonitor,
			),
		)
		if err != nil {
			http.Error(w, "获取 owners 列表失败", http.StatusBadGateway)
			return
		} else {
			defer resp.Body.Close()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	go func() {
		err := http.ListenAndServe(
			net.JoinHostPort(
				"",
				strconv.Itoa(chaindb_config.TCPPortUserService),
			),
			server,
		)
		log.Fatalln("[  User ] Error:", err)
	}()
}
