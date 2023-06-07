INSERT INTO users (social_id, name) VALUES ('karten', 'guest');

INSERT INTO files (id, storage_object_id, name, mime_type, size) VALUES
  ('fd5f451d-fac6-4bc7-a677-34adb39a6701', 'covers_1.jpg', 'Cover', 'image/jpeg', 0),
  ('2f146153-ee2f-4968-a241-11a4f00bf212', 'covers_2.jpg', 'Cover', 'image/jpeg', 0),
  ('40a49f7a-f3b4-4dc4-8b8d-5a7c8e6ccd8b', '40a49f7a-f3b4-4dc4-8b8d-5a7c8e6ccd8b', 'Godot.mp3', 'audio/mpeg', 2399921)
;

INSERT INTO default_cover_images (id) VALUES
  ('fd5f451d-fac6-4bc7-a677-34adb39a6701'),
  ('2f146153-ee2f-4968-a241-11a4f00bf212')
;

INSERT INTO projects (id, user_id, name) VALUES
  ('fd5f451d-fac6-4bc7-a677-34adb39a6701', 1, 'Personal'),
  ('2f146153-ee2f-4968-a241-11a4f00bf212', 1, 'Business'),
  ('2d2712eb-266d-4626-b017-697a67907e28', 1, 'Just Do It vol. 2023'),
  ('1f894df2-f233-4885-81ef-e21aee62e2cd', 1, 'Science pals'),
  ('dd867cb8-3056-49ef-9457-5b7140d16858', 1, 'Other')
;

INSERT INTO boards (
  project_id,
  user_id,
  id,
  name,
  archived,
  color
) VALUES
  ('fd5f451d-fac6-4bc7-a677-34adb39a6701', 1, '29e247c3-69f1-4397-8bab-b1dd10ae28b2', 'Pet projects',  false, 0),
  ('fd5f451d-fac6-4bc7-a677-34adb39a6701', 1, 'f3fc69f2-27aa-4aed-842e-9ed544661bfd', 'Health',        false, 0),
  ('fd5f451d-fac6-4bc7-a677-34adb39a6701', 1, 'ea716cd0-0d2b-4aa9-9a00-e5fce1f6670a', 'TODO Books',    false, 0),
  ('fd5f451d-fac6-4bc7-a677-34adb39a6701', 1, '606ecfd6-2a49-4cc2-911c-a0113ebcf0e6', 'Guitar stuff',  false, 0),
  ('2f146153-ee2f-4968-a241-11a4f00bf212', 1, '311fef19-eb2a-4c04-98d7-3653d2271293', 'Ideas',         false, 0),
  ('2f146153-ee2f-4968-a241-11a4f00bf212', 1, '250bf7a7-ad51-4d62-85cd-c554e4d5f686', 'Householding',  false, 0),
  ('2f146153-ee2f-4968-a241-11a4f00bf212', 1, gen_random_uuid(), 'Marketplace',   false, 0),
  ('2d2712eb-266d-4626-b017-697a67907e28', 1, gen_random_uuid(), 'New Year resolutions', false, 0),
  ('1f894df2-f233-4885-81ef-e21aee62e2cd', 1, gen_random_uuid(), 'Collaboration',        false, 0),
  ('1f894df2-f233-4885-81ef-e21aee62e2cd', 1, gen_random_uuid(), 'Open Source',   false, 0)
;

INSERT INTO task_lists (
  id,
  board_id,
  user_id,
  name,
  position,
  archived,
  color
) VALUES
  ('2fcff999-fa20-419b-84b9-023d81a7688e', '29e247c3-69f1-4397-8bab-b1dd10ae28b2', 1, 'In progress', 100, false, 0),
  ('32f0de22-cc36-4604-9187-f115b45662bd', '29e247c3-69f1-4397-8bab-b1dd10ae28b2', 1, 'Ideas', 300, false, 0),
  ('3b8fcd44-4b59-4fa8-ae12-6ca22ddabd01', '29e247c3-69f1-4397-8bab-b1dd10ae28b2', 1, 'Done', 200, false, 0),
  ('a68f7124-26f0-420a-bf0e-7b4fe27a912e', 'f3fc69f2-27aa-4aed-842e-9ed544661bfd', 1, 'Sport', 100, false, 0),
  ('93892ed8-bd3d-4f8e-b820-bf9a5043bc1d', 'f3fc69f2-27aa-4aed-842e-9ed544661bfd', 1, 'Food', 200, false, 0),
  ('641c10fd-d2c6-44b4-a244-412e28c9edc5', 'ea716cd0-0d2b-4aa9-9a00-e5fce1f6670a', 1, 'Reading', 100, false, 0),
  ('5949e1fe-c8c3-412d-85e4-2adb34bbb1a1', 'ea716cd0-0d2b-4aa9-9a00-e5fce1f6670a', 1, 'Done', 200, false, 0),
  ('7a43965d-048e-4272-9d00-d15f67bbeb81', '606ecfd6-2a49-4cc2-911c-a0113ebcf0e6', 1, 'Good lessons', 100, false, 0),
  (gen_random_uuid(), '250bf7a7-ad51-4d62-85cd-c554e4d5f686', 1, 'Expenses', 100, false, 0)
;

INSERT INTO labels (
  board_id,
  user_id,
  name,
  color
) VALUES
  ('29e247c3-69f1-4397-8bab-b1dd10ae28b2', 1, 'Will not do', 11552306),
  ('29e247c3-69f1-4397-8bab-b1dd10ae28b2', 1, 'Canceled', 11552306)
;

INSERT INTO tasks (
  id,
  task_list_id,
  user_id,
  name,
  text,
  position
) VALUES
  ('3a0c9a3b-bbec-4047-9822-1c4806c2a258', '2fcff999-fa20-419b-84b9-023d81a7688e', 1, 'Refactor geometry', '', 100),
  (gen_random_uuid(), '32f0de22-cc36-4604-9187-f115b45662bd', 1, 'Add tests', '', 100),
  (gen_random_uuid(), '3b8fcd44-4b59-4fa8-ae12-6ca22ddabd01', 1, 'Refactor map loading', '', 100),
  (gen_random_uuid(), '3b8fcd44-4b59-4fa8-ae12-6ca22ddabd01', 1, 'Update sprites', '', 200),
  ('3ea8ea06-d2d2-40af-8c5c-488fc5c9a394', 'a68f7124-26f0-420a-bf0e-7b4fe27a912e', 1, 'Run run run run run run run run run run run run run run run run run run run run run run run run run run run run run run run run run run run run', '', 100),
  (gen_random_uuid(), 'a68f7124-26f0-420a-bf0e-7b4fe27a912e', 1, 'WalkWalkWalkWalkWalkWalkWa lkWalkWalkWal kWalkWalkWalkWalkW alkWalkWalkWalkWalkWalkWalkWalk\nWalkWalkWalkWalkWalkWalkWalkWalkWalkWalkWalkWalkWalkWalk', '', 200),
  ('522a2569-caf5-4c59-8d95-5670ed8378d3', '93892ed8-bd3d-4f8e-b820-bf9a5043bc1d', 1, 'Chips', '', 100),
  ('0f10f18a-bd51-4822-a44f-f5786baf5d07', '93892ed8-bd3d-4f8e-b820-bf9a5043bc1d', 1, 'Noodles', '', 200),
  (gen_random_uuid(), '93892ed8-bd3d-4f8e-b820-bf9a5043bc1d', 1, 'Wok', 'recipe recipe recipe recipe recipe recipe recipe recipereciperecipe recipe', 300),
  ('422eafc4-e488-47b3-9eb4-efd99162bd3d', '641c10fd-d2c6-44b4-a244-412e28c9edc5', 1, 'Little Zaches called Cinnabar', '', 100),
  (gen_random_uuid(), '641c10fd-d2c6-44b4-a244-412e28c9edc5', 1, 'The Prince and the Pauper', '', 200),
  (gen_random_uuid(), '5949e1fe-c8c3-412d-85e4-2adb34bbb1a1', 1, 'The Shadow over Innsmouth', '', 100),
  (gen_random_uuid(), '7a43965d-048e-4272-9d00-d15f67bbeb81', 1, 'https://www.youtube.com/watch?v=Ayx0fEe7gcg', '', 100),
  (gen_random_uuid(), '7a43965d-048e-4272-9d00-d15f67bbeb81', 1, 'https://www.youtube.com/watch?v=Ayx0fEe7gcg', '', 200)
;

INSERT INTO task_labels (label_id, task_id) VALUES
  (1, '3a0c9a3b-bbec-4047-9822-1c4806c2a258'),
  (2, '3a0c9a3b-bbec-4047-9822-1c4806c2a258');

INSERT INTO comments (task_id, user_id, text) VALUES
  ('522a2569-caf5-4c59-8d95-5670ed8378d3', 1, 'https://www.jamieoliver.com/recipes/vegetables-recipes/the-perfect-chips/'),
  ('3a0c9a3b-bbec-4047-9822-1c4806c2a258', 1, 'Read a book about computational geometry'),
  ('3ea8ea06-d2d2-40af-8c5c-488fc5c9a394', 1, '10 km'),
  ('3ea8ea06-d2d2-40af-8c5c-488fc5c9a394', 1, '7 km, trip to a supermaket'),
  ('3ea8ea06-d2d2-40af-8c5c-488fc5c9a394', 1, '0 km, I was lazy af'),
  ('3a0c9a3b-bbec-4047-9822-1c4806c2a258', 1, 'Checkout Coursera'),
  ('422eafc4-e488-47b3-9eb4-efd99162bd3d', 1, 'See the movie'),
  ('0f10f18a-bd51-4822-a44f-f5786baf5d07', 1, 'Buy mint meat')
;

INSERT INTO task_files (task_id, file_id) VALUES
  ('3a0c9a3b-bbec-4047-9822-1c4806c2a258', '40a49f7a-f3b4-4dc4-8b8d-5a7c8e6ccd8b');
