SHELL = /bin/bash -O globstar
CXX ?= c++
CXXFLAGS :=  -Wpedantic -Wall -W -O0 -g3 -Wno-c++20-extensions $(CXXFLAGS)

.PHONY: run
run: | /tmp/shynur/chaindb/user-client/cxx/
	$(CXX) $(CXXFLAGS)  \
		./example.cpp -o /tmp/shynur/chaindb/user-client/cxx/example.exe
	/tmp/shynur/chaindb/user-client/cxx/example.exe

/tmp/shynur/chaindb/user-client/cxx/:
	mkdir -p /tmp/shynur/chaindb/user-client/cxx/
