((auto-mode-alist . (("/.git/COMMIT_EDITMSG\\'" . diff-mode)))
 (nil . ((delete-trailing-whitespace . t)
         (eval . (when (derived-mode-p 'text-mode 'prog-mode 'conf-mode)
                   (add-hook 'before-save-hook
                             #'delete-trailing-whitespace
                             nil "local")))

         (require-final-newline . t)

         (eval . (line-number-mode -1))
         (mode . display-line-numbers)

         (mode . column-number)

         (sentence-end-double-space . t)

         (mode . rainbow)

         (eval . (when (and buffer-file-name (string-match-p "\\.log\\.txt$" buffer-file-name))
                   (auto-revert-tail-mode)))

         (treesit-font-lock-level . 4)))
 (makefile-mode . ((whitespace-style . (face tabs))
                   (mode . whitespace))))
