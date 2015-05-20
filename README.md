
A script which runs analytics on your facebook news feed and pictures to see
who likes your posts the most. Creepy, right?

## Generating a debug token

You need to include a debug token in the golang source code up at the top
under DEBUG_ACCESS_TOKEN

To get this, go to https://developers.facebook.com/tools/explorer/, click
"Get Token" -> "Get Access Token", then check the permission boxes for

* "user_posts"
* "user_likes"

I think that's all you need. Or just check all of them. Whatevs.

Then copy the big string into that field in the go program and run it.
