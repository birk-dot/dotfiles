function ColorMyPencils(color)
	color = color or "vague"
	vim.cmd.colorscheme(color)

	vim.api.nvim_set_hl(0, "Normal", { bg = "none" })
	vim.api.nvim_set_hl(0, "NormalFloat", { bg = "none" })

end

return {
    {
        "folke/tokyonight.nvim",
        config = function()
            require("tokyonight").setup({
                -- your configuration comes here
                -- or leave it empty to use the default settings
                style = "storm", -- The theme comes in three styles, `storm`, `moon`, a darker variant `night` and `day`
                transparent = true, -- Enable this to disable setting the background color
                terminal_colors = true, -- Configure the colors used when opening a `:terminal` in Neovim
                styles = {
                    -- Style to be applied to different syntax groups
                    -- Value is any valid attr-list value for `:help nvim_set_hl`
                    comments = { italic = false },
                    keywords = { italic = false },
                    -- Background styles. Can be "dark", "transparent" or "normal"
                    sidebars = "dark", -- style for sidebars, see below
                    floats = "dark", -- style for floating windows
                },
            })
        end
    },

    {
        "rose-pine/neovim",
        name = "rose-pine",
        config = function()
            require('rose-pine').setup({
                disable_background = true,
                styles = {
                    italic = false,
                },
            })

            vim.cmd("colorscheme rose-pine-moon")

            ColorMyPencils()
        end
    },
    {
        "ellisonleao/gruvbox.nvim",
        name = "gruvbox",
        config = function ()
            require('gruvbox').setup({
                styles = {
                    italic = false,
                },
            })
            vim.cmd("colorscheme gruvbox")

            ColorMyPencils()
        end
    },
    {
      "vague2k/vague.nvim",
      name = "vague",
      lazy = false, -- make sure we load this during startup if it is your main colorscheme
      priority = 1000, -- make sure to load this before all the other plugins
      config = function()
        -- NOTE: you do not need to call setup if you don't want to.
        require("vague").setup({
              transparent = true,
              italic = false,
          -- optional configuration here
            styles = {
                },
        })
        vim.cmd("colorscheme vague")

            ColorMyPencils()
      end
    }


}
