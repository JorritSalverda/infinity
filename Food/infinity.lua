local name = "infinity"
local version = "0.2.19"

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
            sha256 = "9b7a2f79a7ede4f1bbfe203d69bbfb31f682ffdb0e87a802f5a20e5306919d88",
            resources = {
                {
                    path = name .. "-v" .. version .. "-darwin-amd64",
                    installpath = "bin/" .. name,
                    executable = true
                }
            }
        },
        {
            os = "linux",
            arch = "amd64",
            url = "https://github.com/JorritSalverda/" .. name .. "/releases/download/v" .. version .. "/" .. name .. "-v" .. version .. "-linux-amd64.zip",
            -- shasum of the release archive
            sha256 = "4b0df86293982de965ed46ba56f734cb9f870922ea28244e581bf9ff4a86b341",
            resources = {
                {
                    path = name,
                    installpath = "bin/" .. name,
                    executable = true
                }
            }
        },
        {
            os = "windows",
            arch = "amd64",
            url = "https://github.com/JorritSalverda/" .. name .. "/releases/download/v" .. version .. "/" .. name .. "-v" .. version .. "windows-amd64.zip",
            -- shasum of the release archive
            sha256 = "0019dfc4b32d63c1392aa264aed2253c1e0c2fb09216f8e2cc269bbfb8bb49b5",
            resources = {
                {
                    path = name .. ".exe",
                    installpath = "bin\\" .. name .. ".exe"
                }
            }
        }
    }
}
