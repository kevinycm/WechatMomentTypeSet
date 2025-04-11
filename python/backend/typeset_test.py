
from typeset import LayoutEngine

# 测试案例
test_entry = {
    "id": 123,
    "time": "2025-03-20 12:30:15",
    "text": "这是一个需要跨多页的长文本..." * 500,
    "pictures": [
        {
            "width": 1080, 
            "height": 1620, 
            "url": "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkFyklENiaLE4qZQMa4f8rViacHtmAO6RicznsqlO4iaa0LJI"
        },
        {
            "width": 1620,
            "height": 1080,
            "url": "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkUp1vXcxc4J0kVhFwxPRnkLNqvFqER99R2FC0BDmCYx8",
        },
        {
            "width": 810,
            "height": 1080,
            "url": "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnorDmccicxxRQhREKbJjQGuQocfNrSvyvrytnoGHcwSWmDHFQTMbdDj4"
        },
        {
            "width": 810,
            "height": 1080,
            "url": "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9MyDLXiaZ7YtL7DgIDWPqNMS8odq91EdX586jQx2UDvlo"
        },
        {
            "width": 1080,
            "height": 1440,
            "url": "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9rhKE57uzPHBIt4ldv1btOMa0ibW5zxlKRYXQaQMico61Q"
        },
        {
            "width": 810,
            "height": 1080,
            "url": "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9vlFfMRibmJ3m0nJR94FaibxBFEiaw3Dq3UM3fs7cD1ReqA"
        },
        {
            "width": 1080,
            "height": 1620,
            "url": "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFU0ic58W0HarI0kZFBdia9cXEwibbzG2RqEYr0bmiaYMwJ7E6OwMp6haQkk"
        },
        {
            "width": 1080,
            "height": 1440,
            "url": "https://img.diandibianji.com/8u9KefYVGSVTCWiabXdfPI8xRknynoVGZ1rC9M13tRXSU3A1libL6xT8eTkbrtRcRtXOR2C33FSU8"
        },
        {
            "width": 1440,
            "height": 1080,
            "url": "https://img.diandibianji.com/8u9KefYVGSVTCWiabXdfPIib3W7TLFloYBj0w7WOWtxxawey8bHgg4Tyqzrkwre1V8dNA7AlQj4fc"
        }

    ]
}


if __name__ == "__main__":
    engine = LayoutEngine()
    pages = engine.layout_entry(test_entry)
    
    for page in pages:
        print(f"Page {page.page_number}")
        for entry in page.entries:
            print(f"  Entry {entry.id} {'(has time)' if entry.time else ''}")
            print(f"    Text regions: {len(entry.text_regions)}")
            print(f"    Image regions: {len(entry.image_regions)}")
            print(f"    Remaining height: {engine.remaining_height}")