# Section 

```
# Base URL 
https://gall.dcinside.com/board/lists/?id=ohmygirl&page=1
```

## Section 2

```
# Request URL By Page Number Query
https://gall.dcinside.com/board/lists/?id=ohmygirl&page=1
https://gall.dcinside.com/board/lists/?id=ohmygirl&page=2
https://gall.dcinside.com/board/lists/?id=ohmygirl&page=3
https://gall.dcinside.com/board/lists/?id=ohmygirl&page=...
```

## Section 3

```
# Get Post Infomation
1. Post Number
2. Post Title
3. Post Writer
4. Post Write Date
```

## Section 4

```
# Post Struct
post_struct = {
    "post_title": "title",
    "post_writer": "writer",
    "post_write_date": "write_date"
}
# map[post_number] = post_struct
```

## Section 5

```
# Save Post Infomation
Table Info::
    - post_number
    - post_title
    - post_writer
    - post_write_date

# If post_number is already exist in table -> no insert
# If post_number is not exist in table -> insert
```

## Section 6

```
# use Cron Job
1. Get Post Infomation ( witch 1 hour ) ( 1 ~ 100 page )
2. loop for map code
    - 1. check post_number is exist in table
    - 2. if exist -> no insert
    - 3. if not exist -> insert
```
