# WeChat Moment Layout Requirements

## Data Entry Types

There are three types of data entries:

### 1. Time and Text
```json
{
    "id": 121,
    "time": "2025-03-20 12:30:15",
    "text": """
    早上刚出发的时候，
    她接到了一个电话，
    "你出发没？"，
    "刚出发，怎么了？"，
    "哦，没事，我也马上出发。"，
    "好，我8点46到。"，
    "哦，我8点52到（小朋友的妈妈告诉他的）。"，
    两人像老朋友马上要见面一样的聊着……，

    小朋友好像是叫涵涵，
    我以为是个女孩子，
    然后……，
    是个男孩子，

    到了之后才发现，
    今天穿的都是purple，
    可明明这件紫色外套是我去顺丰快递点刚拿的，
    意外，
    又那么巧合，

    第三十七次记录，
    关于你的#CAMPNOW
    """,
    "pictures": []
}
```

### 2. Time and Images (1-9 images)
```json
{
    "id": 121,
    "time": "2025-03-20 12:30:15",
    "text": "",
    "pictures": [
        {
            "width": 200,
            "height": 300,
            "url": "https://img.diandibianji.com/1.jpg"
        },
        {
            "width": 200,
            "height": 300,
            "url": "https://img.diandibianji.com/2.jpg"
        }
    ]
}
```

### 3. Time, Text, and Images (1-9 images)
Same format as above, but with both text and pictures fields populated.

## Algorithm Requirements

1. Layout data entries on A4 paper (3508 x 2480 pixels)
2. Each entry requires its own layout
3. Multi-page handling:
   - 3.1) Short text: Single page
   - 3.2) Long text: Multiple pages
   - 3.3) Images: Pages determined by count
   - 3.4) Text + Images: Combined logic
   - 3.5) Time appears only on first page of multi-page entries

4. Time placement at the top of the first page
5. Text follows time with default width/height for characters
6. Image layout follows 9-grid algorithm:
   - 1 image: [1]
   - 2 images: [2]
   - 3 images: [3]
   - 4 images: [2,2]
   - 5 images: [2,3] or [3,2]
   - 6 images: [3,3]
   - 7 images: [2,2,3], [2,3,2], or [3,2,2]
   - 8 images: [2,3,3], [3,2,3], or [3,3,2]
   - 9 images: [3,3,3]

7. Image handling:
   - Only dimensions provided, no actual images
   - Maintain aspect ratio
   - No cropping allowed
   - A4 size: 3508 x 2480 pixels
   - Row height consistency
   - Row width must fill page (minus margins)
   - Overflow to next page if needed
   - Single image: Fill width if height allows, otherwise prioritize height

## Output Format

```json
[
    {
        "page": 1,
        "time_area": [[start_x, start_y], [end_x, end_y]],
        "time": "2025-03-20 12:12:12",
        "text_area": [[start_x, start_y], [end_x, end_y]],
        "text": "xxxxxxxxx",
        "pictures": [
            {
                "index": 1,
                "area": [[start_x, start_y], [end_x, end_y]],
                "url": "https://img.diandibianji.com/1.jpg"
            }
        ]
    }
]
```

## Implementation Requirements

1. Python implementation preferred
2. Web server using Tornado
3. Single page frontend (index.html)
4. Test cases for all scenarios