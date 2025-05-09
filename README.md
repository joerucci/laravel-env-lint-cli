# laravel-env-lint-cli

A fast, standalone Go CLI to **lint** and **validate** your Laravel-style `.env` files against a YAML schema.

[![Latest Release](https://img.shields.io/github/v/release/joerucci/laravel-env-lint-cli?style=flat-square)](https://github.com/joerucci/laravel-env-lint-cli/releases)
![CI](https://github.com/joerucci/laravel-env-lint-cli/actions/workflows/ci.yml/badge.svg)
[![Go Version](https://img.shields.io/badge/go-1.24.3-blue?style=flat-square)](https://golang.org/doc/go1.24)
[![License](https://img.shields.io/github/license/joerucci/laravel-env-lint-cli?style=flat-square)](LICENSE)

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