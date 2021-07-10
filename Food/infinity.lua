local name = "infinity"
local version = "0.2.15"

food = {
    name = name,
    description = "Infinity is a CLI to easily build your applications using a pipeline as code",
    homepage = "https://github.com/JorritSalverda/infinity",
    version = version,
    packages = {
        {
            os = "darwin",
            arch = "amd64",
            url = "https://github.com/JorritSalverda/" .. name .. "/releases/download/v" .. version .. "/" .. name .. "-v" .. version .. "-darwin-amd64.zip",
            -- shasum of the release archive
            sha256 = "b210399a32c0189bdc1c46571cb19ef09b3843952eb9f089fa2c61c1c6a9f7ab",
            resources = {
                {
                    path = name .. "-v" .. version .. "-darwin-amd64",
                    installpath = "bin/" .. name,
                    executable = true
                }
            }
        }-- ,
        -- {
        --     os = "linux",
        --     arch = "amd64",
        --     url = "https://github.com/JorritSalverda/" .. name .. "/releases/download/v" .. version .. "/" .. name .. "-v" .. version .. "-linux-amd64.zip",
        --     -- shasum of the release archive
        --     sha256 = "",
        --     resources = {
        --         {
        --             path = name,
        --             installpath = "bin/" .. name,
        --             executable = true
        --         }
        --     }
        -- },
        -- {
        --     os = "windows",
        --     arch = "amd64",
        --     url = "https://github.com/JorritSalverda/" .. name .. "/releases/download/v" .. version .. "/" .. name .. "-v" .. version .. "windows-amd64.zip",
        --     -- shasum of the release archive
        --     sha256 = "",
        --     resources = {
        --         {
        --             path = name .. ".exe",
        --             installpath = "bin\\" .. name .. ".exe"
        --         }
        --     }
        -- }
    }
}
