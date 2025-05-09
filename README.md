# laravel-env-lint-cli

A fast, standalone Go CLI to **lint** and **validate** your Laravel-style `.env` files against a YAML schema.

![CI](https://github.com/joerucci/laravel-env-lint-cli/actions/workflows/ci.yml/badge.svg)

## Features

- ✅ **Validate** environment variables types: `string`, `integer`, `boolean`, `float`  
- 🔒 **Decrypt** and validate encrypted `.env.encrypted` files (AES-256-CBC + PKCS#7)  
- 🚀 **Laravel-style** casting of `true`/`false`/`null`/`empty`  
- 🛡️ Enforce `required`, `nullable`, `one_of`, and conditional `required_when` rules  
- 🧪 Fully tested core logic and loader

## Installation

### From source

```bash
go install github.com/joerucci/laravel-env-lint-cli/cmd/laravel-env-lint@latest