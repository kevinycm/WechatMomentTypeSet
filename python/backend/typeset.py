import math

class LayoutEngine:
    def __init__(self, entry):
        self.entry = entry
        self.pages = []
        self.current_page = None
        self.margin_left = 100
        self.margin_right = 100
        self.margin_top = 100
        self.margin_bottom = 100
        self.available_width = 2480 - self.margin_left - self.margin_right
        self.available_height = 3508 - self.margin_top - self.margin_bottom
        self.time_height = 100
        self.line_height = 40  # 假设每行文本高度为20像素

    def process_entry(self):
        self.new_page()
        self.add_time()
        self.process_text(self.entry.get('text', ''))
        self.process_pictures(self.entry.get('pictures', []))
        return self.pages

    def new_page(self):
        self.current_page = {
            'page': len(self.pages) + 1,
            'time_area': None,
            'time': None,
            'text_areas': [],
            'texts': [],
            'pictures': []
        }
        self.pages.append(self.current_page)
        self.current_y = self.margin_top

    def add_time(self):
        if self.current_page['page'] != 1:
            return
        time_str = self.entry['time']
        x0 = self.margin_left
        y0 = self.margin_top
        x1 = x0 + self.available_width
        y1 = y0 + self.time_height
        self.current_page['time_area'] = [(x0, y0), (x1, y1)]
        self.current_page['time'] = time_str
        self.current_y = y1

    def process_text(self, text):
        if not text.strip():
            return
        lines = text.split('\n')
        current_line = 0
        while current_line < len(lines):
            if self.current_page['page'] == 1:
                remaining_height = self.available_height - self.time_height
            else:
                remaining_height = self.available_height
            available_lines = min(len(lines) - current_line, 
                                 remaining_height // self.line_height)
            chunk = lines[current_line:current_line + available_lines]
            self.add_text_chunk(chunk)
            current_line += available_lines
            if current_line < len(lines):
                self.new_page()

    def add_text_chunk(self, chunk):
        start_y = self.current_y
        text_height = len(chunk) * self.line_height
        area = [
            (self.margin_left, start_y),
            (self.margin_left + self.available_width, start_y + text_height)
        ]
        self.current_page['text_areas'].append(area)
        self.current_page['texts'].append('\n'.join(chunk))
        self.current_y += text_height

    def process_pictures(self, pictures):
        layout = self.get_layout(len(pictures))
        current_idx = 0
        for row in layout:
            row_pics = pictures[current_idx:current_idx+row]
            current_idx += row
            self.process_picture_row(row_pics)

    def get_layout(self, n):
        layout_rules = {
            1: [1], 2: [2], 3: [3], 4: [2,2],
            5: [2,3], 6: [3,3], 7: [2,2,3],
            8: [3,3,2], 9: [3,3,3]
        }
        return layout_rules.get(n, [])

    def process_picture_row(self, row_pics):
        while True:
            available_height = self.available_height - (self.current_y - self.margin_top)
            if available_height <= 0:
                self.new_page()
                continue
            
            total_width = 0
            common_height = available_height

            scaled_widths = []
            for pic in row_pics:
                aspect_ratio = pic['width'] / pic['height']
                scaled_width = common_height * aspect_ratio
                scaled_widths.append(scaled_width)
                total_width += scaled_width

            if len(row_pics) == 1:
                if common_height <= available_height:
                    self.place_pictures(row_pics, scaled_widths, common_height)
                    break
                else:
                    self.new_page()
                    continue
            else:
                # 如果总宽度超过可用宽度，则需要调整
                if total_width > self.available_width:
                    # 计算缩放因子
                    scale_factor = self.available_width / total_width
                    common_height *= scale_factor
                    total_width = self.available_width
                    scaled_widths = [w * scale_factor for w in scaled_widths]
                
                # 检查调整后的高度是否仍然适合
                if common_height <= available_height:
                    if abs(total_width - self.available_width) > 1:
                        self.new_page()
                        continue

                    self.place_pictures(row_pics, scaled_widths, common_height)
                    break
                else:
                    self.new_page()

    def place_pictures(self, row_pics, scaled_widths, common_height):
        start_y = self.current_y
        x = self.margin_left
        for i, (pic, pic_width) in enumerate(zip(row_pics, scaled_widths)):
            area = [
                (x, start_y),
                (x + pic_width, start_y + common_height)
            ]
            self.current_page['pictures'].append({
                'index': i+1,
                'area': area,
                'url': pic['url']
            })
            x += pic_width
        self.current_y += common_height