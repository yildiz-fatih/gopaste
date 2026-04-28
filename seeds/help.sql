INSERT INTO pastes (slug, content, created, expires)
VALUES ('help', $$
# welcome to the world of gopaste!

> our servers are replaceable, your paste isn't

## basic usage

1. type (or paste) code
2. press 'ctrl + s' to save
3. share URL

## why i built this

every other pastebin felt like it was frozen in 1999!

so i thought:

> yeah, well...  
> i'm gonna go build my own pastebin,  
> with syntax highlighting and keyboard navigation
>
> in fact, forget the pastebin!

$$, NOW(), NOW() + INTERVAL '100 years')
ON CONFLICT (slug) DO UPDATE SET content = EXCLUDED.content;
