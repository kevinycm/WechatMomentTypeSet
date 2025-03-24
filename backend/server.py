import tornado.ioloop
import tornado.web
import json

from typeset import LayoutEngine


sample_data = {
    121: {
        "id": 121,
        "time": "2025-03-20 12:30:15",
        "text": "这是一个需要跨多页的长文本...\n" * 24,
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
    },
    122: {
        "id": 122,
        "time": "2025-03-20 12:30:15",
        "text": "这是一个需要跨多页的长文本...\n" * 24,
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
            }
        ]
    },
    123: {
        "id": 123,
        "time": "2025-03-20 12:30:15",
        "text": "这是一个需要跨多页的长文本...\n" * 24,
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
            }
        ]
    },
    124: {
        "id": 124,
        "time": "2025-03-20 12:30:15",
        "text": "这是一个需要跨多页的长文本...\n" * 24,
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
            }
        ]
    },
    125: {
        "id": 125,
        "time": "2025-03-20 12:30:15",
        "text": "这是一个需要跨多页的长文本...\n" * 24,
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
            }
        ]
    },
    126: {
        "id": 126,
        "time": "2025-03-20 12:30:15",
        "text": "这是一个需要跨多页的长文本...\n" * 24,
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
            }
        ]
    },
    127: {
        "id": 127,
        "time": "2025-03-20 12:30:15",
        "text": "这是一个需要跨多页的长文本...\n" * 24,
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
            }
        ]
    },
    128: {
        "id": 128,
        "time": "2025-03-20 12:30:15",
        "text": "这是一个需要跨多页的长文本...\n" * 24,
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
            }
        ]
    },
    129: {
        "id": 129,
        "time": "2025-03-20 12:30:15",
        "text": "这是一个需要跨多页的长文本...\n" * 24,
        "pictures": [
            {
                "width": 1080, 
                "height": 1620, 
                "url": "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkFyklENiaLE4qZQMa4f8rViacHtmAO6RicznsqlO4iaa0LJI"
            }
        ]
    },
    130: {
        "id": 130,
        "time": "2025-03-20 12:30:15",
        "text": "这是一个需要跨多页的长文本...\n" * 24,
        "pictures": [
            
        ]
    },
    131: {
        "id": 121,
        "time": "2025-03-20 12:30:15",
        "text": "",
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
}



class MainHandler(tornado.web.RequestHandler):
    def get(self):
        self.render("index.html")

class LayoutHandler(tornado.web.RequestHandler):
    def get(self, entry_id):
        entry = sample_data.get(int(entry_id))
        if not entry:
            self.set_status(404)
            return
        engine = LayoutEngine(entry)
        result = engine.process_entry()
        self.write(json.dumps(result))

def make_app():
    return tornado.web.Application([
        (r"/", MainHandler),
        (r"/layout/(\d+)", LayoutHandler),
    ])

if __name__ == "__main__":
    app = make_app()
    app.listen(8888)
    tornado.ioloop.IOLoop.current().start()