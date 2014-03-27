Find a [panagram][pan], but with two differences:

1. Use each letter *exactly* once, instead of *at least* once.

2. No grammar restrictions -- it doesn't need to be a sentence.

[pan]: http://en.wikipedia.org/wiki/Panagram

## Motivation

I was playing with my infant son's blocks. He has 26 of them, one for each
letter. I got to wondering how many words I could spell with those blocks.
Could I use them all?

## Thoughts on dictionary

After looking at some perfect pangrams online it became clear that I didn't
want to use a full dictionary. [For example][w]:

    Cwm fjord veg balks nth pyx quiz.

[w]: http://en.wikipedia.org/wiki/List_of_pangrams

Sure it's a perfect pangram, but it won't impress anybody because it looks like
a monkey typed it out.

My solution was to restrict my search to commonly used words.

## Results

My best results after using the 10,000 most commonly used words found in all
[Project Gutenberg][pg] books:

[pg]: http://www.gutenberg.org/

    24 [dry] [length] [fix] [backs] [jump] [vow]
    24 [things nights] [verb] [wax] [dock] [fly] [jump]
    24 [why] [bricks] [fox] [vent] [glad] [jump]
    24 [why] [drag] [blocks] [fix] [vent] [jump]
    24 [why] [grand] [block] [vest] [fix] [jump]
    24 [why] [grind] [fox] [vest] [black] [jump]

Note that none use either a `q` or a `z`.

A search of the full dictionary yielded some interesting results. When
restricting my search to words that contain both a `q` and `z` I found no
combinations that used more than 24 letters (which means that no combinations
use all 26 letters). Same when I required only a `q`.  *But*, when I required
only a `z` the program took off. Most runs prior to that took 0.5-2 hours.
After *18* hours I ended the `z`-only run with *1200* combinations banked up
that used 25 unique letters.

Moral of the story is that `q` is a much less friendly letter than `z`.

Here's a sample 25-long solution. The rest (before I gave my poor laptop a
break) can be found in `partial-z-only.txt`.

    25 [dwarf] [thumb] [jock] [vex] [zings] [ply]
