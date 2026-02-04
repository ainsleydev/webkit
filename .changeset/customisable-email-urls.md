---
"@ainsleydev/payload-helper": minor
---

Adding customisable URL callbacks for verification and forgot password emails. Users can now provide a `url` callback in their email configuration to generate custom URLs for frontend verification flows.

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
