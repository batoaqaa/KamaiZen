# KamaiZen
A language server for kamailio configuration files

For syntax highlighting, use [tree-sitter-kamailio-cfg](https://github.com/IbrahimShahzad/tree-sitter-kamailio-cfg)


## Features
- [ ] Code completion
    - [ ] Variables
    - [x] Functions
    - [ ] Modules
    - [ ] Keywords
    - [ ] Parameters
    - [x] Headers
- [ ] Code navigation
  - [ ] Go to definition - In progress
  - [ ] Find references - In progress
- [ ] Code formatting
- [ ] Code folding
- [ ] Diagnostics
    - [x] Syntax Errors
    - [x] Invalid statements
    - [x] Unreachable code
    - [ ] Unused variables
    - [ ] Unused functions
    - [ ] Unused modules
    - [ ] Unused parameters
- [ ] Hover
    - [x] Show documentation
    - [ ] Show type
- [ ] Highlight


## Installation

### From source

```bash
git clone https://github.com/IbrahimShahzad/KamaiZen.git 
cd KamaiZen
go build
cp KamaiZen /usr/local/bin
```

### To use with Neovim

update your `init.lua` with the following

```lua
local client = vim.lsp.start_client {
  cmd = { '/usr/local/bin/KamaiZen' }, -- Path to KamaiZen executable
  name = 'KamaiZen',
  settings = {
    kamaizen = {
      logLevel = 1,
      kamailioSourcePath = '/path/to/kamailio-source', -- Path to kamailio source
    },
  },
}

if not client then
  vim.notify('Failed to start LSP client', vim.log.levels.ERROR)
  return
end

vim.api.nvim_create_autocmd('FileType', {
  pattern = 'kamailio',
  callback = function()
    vim.lsp.buf_attach_client(0, client)
  end,
})
```

## Integration

### Neovim

- [x] [kamaizen.nvim](https://github.com/IbrahimShahzad/kamaizen.nvim)

### Vscode

> Not yet available

- [ ] [vscode-kamaizen](github.com/IbrahimShahzad/vscode-kamaizen)


