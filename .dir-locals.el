((auto-mode-alist . (("/.git/COMMIT_EDITMSG\\'" . diff-mode)))
 (nil . ((delete-trailing-whitespace . t)
         (eval . (when (derived-mode-p 'text-mode 'prog-mode 'conf-mode)
                   (add-hook 'before-save-hook
                             #'delete-trailing-whitespace
                             nil "local")))

         (delete-trailing-lines . t)
         (require-final-newline . t)

         (eval . (when buffer-file-name
                   (when (string-match-p "\\`\\(LICENSE\\|License\\|license\\|COPYING\\)$"
                                         (file-name-base buffer-file-name))
                     (setq-local buffer-read-only t))))
         (eval . (when buffer-file-name
                   (when (string-match-p "\\.log\\.txt$"
                                         (file-name-nondirectory buffer-file-name))
                     (setq-local buffer-read-only t))))

         (eval . (line-number-mode -1))
         (mode . display-line-numbers)

         (mode . column-number)

         (project-vc-merge-submodules . nil)

         (sentence-end-double-space . t)

         (auto-revert-verbose . nil)
         (auto-revert-avoid-polling . t)
         (eval . (when buffer-file-name
                   (when (string-match-p "\\.log\\.txt$" buffer-file-name)
                     (auto-revert-tail-mode))))

         (treesit-font-lock-level . 4)))
 (makefile-mode . ((whitespace-style . (face tabs))
                   (mode . whitespace)))
 (yaml-mode . ((indent-tabs-mode . nil)
               (tab-width . 2))))
