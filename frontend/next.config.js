/** @type {import('next').NextConfig} */
const nextConfig = {
    reactStrictMode: true,
}
const withTM = require('next-transpile-modules')(['react-daisyui'])
module.exports = withTM({nextConfig})
