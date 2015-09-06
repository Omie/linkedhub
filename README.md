linkedhub
=========

Contact discovery like LinkedIn, something to learn Go language.

Crawl through user's repositories to find other collaborators and repeat.

Freezing code at this point.
Multiple reasons : 
- Its expensive. Once I start going for depth more than 3, I will exhaust my API limit
before finishing, still that data won't be able to help much.
- It makes sense if Github does it, doesn't make sense if I try to replicate their 
database and find contact relationships.
- This was meant as an exercise to learn Go language, and certainly has helped me a lot
- Time to build something new :)

How to use
----------------
- Make sure you have Go environment set
- Checkout code; cd linkedhub
- export GOPATH=`pwd`
- go build github.com/omie/linkedhub
- ./linkedhub


Demo
-----

Screenshot

![linkedhub](https://raw.githubusercontent.com/Omie/linkedhub/screenshot/linkedhub.png)


License
--------
MIT License

