# Changelog

All notable changes to this project will be documented in this file. See [commit-and-tag-version](https://github.com/absolute-version/commit-and-tag-version) for commit guidelines.

## 1.0.0 (2023-11-29)


### Features

* active role has permission middleware ([04e703f](https://bitbucket.org/dptsi/base-go/commit/04e703f8680ce63a1040beb0c0bc1772b14f8d8a))
* **auth:** handle forgot to add middleware ([1c30c6b](https://bitbucket.org/dptsi/base-go/commit/1c30c6b3d075b85589ef05177387f52a67395b60))
* **auth:** login using oidc ([cb61cbd](https://bitbucket.org/dptsi/base-go/commit/cb61cbda77eadfe4a1aeaa8f8cef62c46d361424))
* **auth:** make debugging easier from response ([367e95f](https://bitbucket.org/dptsi/base-go/commit/367e95fbe122e7f2f8e8b5e56b686003be1f3e74))
* **auth:** redirect to frontend after login ([374e81f](https://bitbucket.org/dptsi/base-go/commit/374e81fe054fe98b2b9a25724767385dd40cfa4f))
* **auth:** response success logged in ([05da084](https://bitbucket.org/dptsi/base-go/commit/05da0841542a95469f76a7d4982bfe4bfc437780))
* **auth:** switch active role route ([354c094](https://bitbucket.org/dptsi/base-go/commit/354c094949c2682807baf81ae4bc6a00c10cedec))
* **auth:** validate nonce ([039a7e3](https://bitbucket.org/dptsi/base-go/commit/039a7e3a8e41695508dd181c39dc345a88824880))
* entra and sso oidc ([637a3a8](https://bitbucket.org/dptsi/base-go/commit/637a3a8860fbad5243ea02b4942a070299db88f3))
* **frs:** record not found error ([4510ec0](https://bitbucket.org/dptsi/base-go/commit/4510ec080c5132068250c423ec2a7e68daa347b5))
* mkdocs endpoint ([f6016b4](https://bitbucket.org/dptsi/base-go/commit/f6016b4467df05b70f979bcd3ad4059035a64862))
* **oidc:** pkce support ([8f2f629](https://bitbucket.org/dptsi/base-go/commit/8f2f629995fc79bf19259ee38d48bac30f85af17))
* permissions from myits sso resource ([fea22b8](https://bitbucket.org/dptsi/base-go/commit/fea22b8f9c9f9f4d7334f7172c11285d48be7dc3))
* permissions from myits sso resource ([6511c05](https://bitbucket.org/dptsi/base-go/commit/6511c05603704de58cc755dc043710b0213d14a6))
* return null instead of empty string ([3815872](https://bitbucket.org/dptsi/base-go/commit/3815872a5a74a83fc5942aefc0f91e9073a5b53a))
* **script:** database setup boilerplate per module ([8bffd4d](https://bitbucket.org/dptsi/base-go/commit/8bffd4d3b9e9d1fe7176c109d2154a9f997eb43c))
* **session:** sql server session driver ([c3924cb](https://bitbucket.org/dptsi/base-go/commit/c3924cb7cc0f7bf2823814f7989a02160b1c0fe0))
* swagger ui endpoint ([ceba556](https://bitbucket.org/dptsi/base-go/commit/ceba556af40c7cf7f8e3e74f314db1f0dd78404c))


### Bug Fixes

* **auth:** empty active role string if it's null ([adb40a1](https://bitbucket.org/dptsi/base-go/commit/adb40a1de31893ded26a34059d815e4cf039c378))
* case sensitive header csrf ([a0859fc](https://bitbucket.org/dptsi/base-go/commit/a0859fc9f526964892df1869eba317d3736b84e4))
* error myits sso ([da63bc1](https://bitbucket.org/dptsi/base-go/commit/da63bc19bcbefd90b1240b804e28155b3d384632))
* firestore still return expired session ([687bee3](https://bitbucket.org/dptsi/base-go/commit/687bee3bc7bd35ff0de55a83bacc347a2776b52d))
* route bug ([b32c9ad](https://bitbucket.org/dptsi/base-go/commit/b32c9ad315e35129675bf85c9ddf9d8b547e08eb))
* **script:** wrong import ([124506b](https://bitbucket.org/dptsi/base-go/commit/124506b64dd2ed680e0dddcc3f39536fdadaa033))
* **script:** wrong import in routes ([3c320f9](https://bitbucket.org/dptsi/base-go/commit/3c320f9cef61e3b49bb13ab0d6a3bd675e1a3f38))
* swagger can't call the api ([153ff82](https://bitbucket.org/dptsi/base-go/commit/153ff82bf01172ebc65176ac625030f8c9aeaa6d))
* swagger errror ([6de3ee2](https://bitbucket.org/dptsi/base-go/commit/6de3ee2942126c634e36a5fcef9a6cd6d05fdeee))
* user doesn;t have active role by default ([8cd1110](https://bitbucket.org/dptsi/base-go/commit/8cd11106f27984729b4c170568ffecd614b82e2e))
* user roles null ([f17fc86](https://bitbucket.org/dptsi/base-go/commit/f17fc86b4cf30c8ca6bd6726fcd119b3292ae6d6))
* wrong http status code ([c50cf3e](https://bitbucket.org/dptsi/base-go/commit/c50cf3ef58bf7c8d15789d39d24f3755f484bb7c))
* wrong role being serialized ([562ebc9](https://bitbucket.org/dptsi/base-go/commit/562ebc9e9cb6d7d08f9deb15fedd5bbb28bb15ab))

## [0.0.1](https://bitbucket.org/dptsi/base-go/compare/v0.0.0...v0.0.1) (2023-11-29)


### Features

* active role has permission middleware ([04e703f](https://bitbucket.org/dptsi/base-go/commit/04e703f8680ce63a1040beb0c0bc1772b14f8d8a))
* **auth:** make debugging easier from response ([367e95f](https://bitbucket.org/dptsi/base-go/commit/367e95fbe122e7f2f8e8b5e56b686003be1f3e74))
* **auth:** switch active role route ([354c094](https://bitbucket.org/dptsi/base-go/commit/354c094949c2682807baf81ae4bc6a00c10cedec))
* entra and sso oidc ([637a3a8](https://bitbucket.org/dptsi/base-go/commit/637a3a8860fbad5243ea02b4942a070299db88f3))
* **frs:** record not found error ([4510ec0](https://bitbucket.org/dptsi/base-go/commit/4510ec080c5132068250c423ec2a7e68daa347b5))
* permissions from myits sso resource ([fea22b8](https://bitbucket.org/dptsi/base-go/commit/fea22b8f9c9f9f4d7334f7172c11285d48be7dc3))
* permissions from myits sso resource ([6511c05](https://bitbucket.org/dptsi/base-go/commit/6511c05603704de58cc755dc043710b0213d14a6))
* return null instead of empty string ([3815872](https://bitbucket.org/dptsi/base-go/commit/3815872a5a74a83fc5942aefc0f91e9073a5b53a))


### Bug Fixes

* error myits sso ([da63bc1](https://bitbucket.org/dptsi/base-go/commit/da63bc19bcbefd90b1240b804e28155b3d384632))
* swagger can't call the api ([153ff82](https://bitbucket.org/dptsi/base-go/commit/153ff82bf01172ebc65176ac625030f8c9aeaa6d))
* user doesn;t have active role by default ([8cd1110](https://bitbucket.org/dptsi/base-go/commit/8cd11106f27984729b4c170568ffecd614b82e2e))
* wrong http status code ([c50cf3e](https://bitbucket.org/dptsi/base-go/commit/c50cf3ef58bf7c8d15789d39d24f3755f484bb7c))
