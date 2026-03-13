# @ainsleydev/payload-helper

## 0.3.1

### Patch Changes

- d57ed1b: Fix slug field component import path resolving to relative path in consumer's importMap

## 0.3.0

### Minor Changes

- fe26c53: Add a reusable `SlugField` helper with a lockable admin component and automatic slug formatting hook.

## 0.2.2

### Patch Changes

- 9eb5cc1: Moving lexical and @lexical/\* to peerDependencies with a >=0.35.0 range, aligning with @payloadcms/richtext-lexical@3.74.0 and allowing consuming projects to inherit whatever version their Payload install resolves.

## 0.2.1

### Patch Changes

- 4bc9eec: Add PublishedAt and URLField field helpers and expose via ./fields export path

## 0.2.0

### Minor Changes

- fd80a0c: Adding customisable URL callbacks for verification and forgot password emails. Users can now provide a `url` callback in their email configuration to generate custom URLs for frontend verification flows.

  Example usage:

  ```ts
  email: {
    forgotPassword: {
      url: ({ token }) => `https://myapp.com/reset?token=${token}`,
    },
    verifyAccount: {
      url: ({ token, collection }) => `https://myapp.com/verify?token=${token}`,
    },
  }
  ```

## 0.1.6

### Patch Changes

- 3233efa: Fix admin Icon and Logo visibility by removing block-level figure wrapper

## 0.1.5

### Patch Changes

- 6906b17: fix: Adding .js extension to ESM import for Node.js compatibility

## 0.1.4

### Patch Changes

- 36eb3da: Decouple email preview CLI from Payload config to fix path alias resolution errors.

## 0.1.3

### Patch Changes

- 20b2aa5: Bumping Deps

## 0.1.2

### Patch Changes

- Updated dependencies [3ea3383]
  - @ainsleydev/email-templates@0.0.4

## 0.1.1

### Patch Changes

- 7f3a0a5: Barrel Exports

## 0.1.0

### Minor Changes

- 77bc7e9: Email Integration

## 0.0.40

### Patch Changes

- bf60a66: Icon Support

## 0.0.39

### Patch Changes

- 7b5f55c: Admin Logo

## 0.0.38

### Patch Changes

- 89f9d22: Fixing GitHub releases

## 0.0.37

### Patch Changes

- e7eb847: Testing pipe

## 0.0.36

### Patch Changes

- 9491312: Testing changeset pipeline
