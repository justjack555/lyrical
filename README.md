# Lyrical

## Summary
You remember a few lyrics to a song you heard earlier in
the day. You type them into your favorite music player - no results.
Sure, Google search will likely yield your target song, but
perhaps the lyrics you remember are a set of "stopwords" that
yield too many results in the Google search context. 

Here's where lyrical steps in. As an extension to your
favorite music player, lyrical will provide a fast mapping
to your target song by querying just the lyrics you know.

## Components
### Parser
Lyrical contains a song parser that, given a repository of song lyric files,
parses them and builds an inverted index that allows for
the fast lookup speeds on a massive set of songs.
 
## Run Lyrical
* Run the song parser: 

    `go run cmd/parser/parser.go cmd/parser/parser.yaml`
    
## Enhancements
### Parser
#### File Loading
* During song lyrics file traversing, we change directory
to the parent prior to each call of processFileInfo on a
child file. This should be improved by either goroutine
usage to create a new process stack or by improved
usage of the Chdir function