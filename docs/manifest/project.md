# Project

Defines project identity and source repository. This object is used for outputting metadata and prefixing resources.

## Attributes

| Key           | Description                        | Required | Notes      |
|---------------|------------------------------------|----------|------------|
| `name`        | Project machine-readable name      | Yes      | kebab-case |
| `title`       | Human-readable project title       | Yes      | Title Case |
| `description` | Description of the project         | Yes      |            |
| `repo`        | Github repository link (HTTPs URL) | Yes      |            |


## Example

```json
{
    "project": {
        "name": "my-website",
        "title": "My Website",
        "description": "My website is a bespoke sales platform for developers and designers.",
        "repo": "git@github.com:ainsley/my-website.git"
    }
}
```