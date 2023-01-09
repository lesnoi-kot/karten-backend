SET search_path TO karten;

BEGIN;

INSERT INTO projects (id, name) VALUES
  ('fd5f451d-fac6-4bc7-a677-34adb39a6701', 'Personal'),
  ('2f146153-ee2f-4968-a241-11a4f00bf212', 'Business'),
  ('2d2712eb-266d-4626-b017-697a67907e28', 'Just Do It vol. 2023'),
  ('1f894df2-f233-4885-81ef-e21aee62e2cd', 'Science pals')
;

INSERT INTO boards (
  project_id,
  id,
  name,
  archived,
  color,
  date_created,
  date_last_viewed,
  cover_url
) VALUES
  ('fd5f451d-fac6-4bc7-a677-34adb39a6701', '29e247c3-69f1-4397-8bab-b1dd10ae28b2', 'Pet projects',  false, 0, '1970-01-01T00:00:00Z', '1970-01-01T00:00:00Z', NULL),
  ('fd5f451d-fac6-4bc7-a677-34adb39a6701', 'f3fc69f2-27aa-4aed-842e-9ed544661bfd', 'Health',        false, 1, '1970-01-01T00:00:00Z', '1970-01-01T00:00:00Z', NULL),
  ('fd5f451d-fac6-4bc7-a677-34adb39a6701', 'ea716cd0-0d2b-4aa9-9a00-e5fce1f6670a', 'TODO Books',    false, 0, '1970-01-01T00:00:00Z', '1970-01-01T00:00:00Z', NULL),
  ('fd5f451d-fac6-4bc7-a677-34adb39a6701', '606ecfd6-2a49-4cc2-911c-a0113ebcf0e6', 'Guitar stuff',  false, 0, '1970-01-01T00:00:00Z', '1970-01-01T00:00:00Z', NULL),
  ('2f146153-ee2f-4968-a241-11a4f00bf212', '311fef19-eb2a-4c04-98d7-3653d2271293', 'Ideas',         false, 0, '1970-01-01T00:00:00Z', '1970-01-01T00:00:00Z', NULL),
  ('2f146153-ee2f-4968-a241-11a4f00bf212', '250bf7a7-ad51-4d62-85cd-c554e4d5f686', 'Householding',  false, 2, '1970-01-01T00:00:00Z', '1970-01-01T00:00:00Z', NULL)
;

INSERT INTO task_lists (
  id,
  board_id,
  name,
  position,
  archived,
  date_created,
  color
) VALUES
  ('2fcff999-fa20-419b-84b9-023d81a7688e', '29e247c3-69f1-4397-8bab-b1dd10ae28b2', 'In progress', 100, false, '1970-01-01T00:00:00Z', 0),
  ('32f0de22-cc36-4604-9187-f115b45662bd', '29e247c3-69f1-4397-8bab-b1dd10ae28b2', 'Ideas', 300, false, '1970-01-01T00:00:00Z', 0),
  ('3b8fcd44-4b59-4fa8-ae12-6ca22ddabd01', '29e247c3-69f1-4397-8bab-b1dd10ae28b2', 'Done', 200, false, '1970-01-01T00:00:00Z', 0),
  ('a68f7124-26f0-420a-bf0e-7b4fe27a912e', 'f3fc69f2-27aa-4aed-842e-9ed544661bfd', 'Sport', 100, false, '1970-01-01T00:00:00Z', 0),
  ('93892ed8-bd3d-4f8e-b820-bf9a5043bc1d', 'f3fc69f2-27aa-4aed-842e-9ed544661bfd', 'Food', 200, false, '1970-01-01T00:00:00Z', 0),
  ('641c10fd-d2c6-44b4-a244-412e28c9edc5', 'ea716cd0-0d2b-4aa9-9a00-e5fce1f6670a', 'Reading', 100, false, '1970-01-01T00:00:00Z', 0),
  ('5949e1fe-c8c3-412d-85e4-2adb34bbb1a1', 'ea716cd0-0d2b-4aa9-9a00-e5fce1f6670a', 'Done', 200, false, '1970-01-01T00:00:00Z', 0),
  ('7a43965d-048e-4272-9d00-d15f67bbeb81', '606ecfd6-2a49-4cc2-911c-a0113ebcf0e6', 'Good lessons', 100, false, '1970-01-01T00:00:00Z', 0)
;

INSERT INTO tasks (
  id,
  task_list_id,
  name,
  text,
  position,
  due_date,
  date_created
) VALUES
  ('3a0c9a3b-bbec-4047-9822-1c4806c2a258', '2fcff999-fa20-419b-84b9-023d81a7688e', 'Refactor geometry', '', 100, NULL, '1970-01-01T00:00:00Z'),
  ('3ea8ea06-d2d2-40af-8c5c-488fc5c9a394', 'a68f7124-26f0-420a-bf0e-7b4fe27a912e', 'Run', '', 100, NULL, '1970-01-01T00:00:00Z'),
  ('522a2569-caf5-4c59-8d95-5670ed8378d3', '93892ed8-bd3d-4f8e-b820-bf9a5043bc1d', 'Chips', 'Text', 100, '1970-01-01T00:00:00Z', '1970-01-01T00:00:00Z'),
  ('0f10f18a-bd51-4822-a44f-f5786baf5d07', '93892ed8-bd3d-4f8e-b820-bf9a5043bc1d', 'Noodles', '', 200, NULL, '1970-01-01T00:00:00Z'),
  ('422eafc4-e488-47b3-9eb4-efd99162bd3d', '641c10fd-d2c6-44b4-a244-412e28c9edc5', 'Little Zaches called Cinnabar', '', 100, NULL, '1970-01-01T00:00:00Z')
;

INSERT INTO comments (id, task_id, author, text, date_created) VALUES
  ('4d715efa-8e2f-4488-99d2-5b69e7a43aec', '522a2569-caf5-4c59-8d95-5670ed8378d3', 'User', 'https://www.jamieoliver.com/recipes/vegetables-recipes/the-perfect-chips/', '1970-01-01T00:00:00Z'),
  ('56b3893f-1887-4b7b-bae9-31f959ecca68', '3a0c9a3b-bbec-4047-9822-1c4806c2a258', 'User', 'Read a book about computational geometry', '1970-01-01T00:00:00Z'),
  ('bb631152-5e05-4e99-b67e-5bd8a1167729', '3ea8ea06-d2d2-40af-8c5c-488fc5c9a394', 'User', '10 km', '1970-01-01T00:00:00Z'),
  ('0dca54e0-66a9-42ca-bb9f-dbf926e2c8e5', '3ea8ea06-d2d2-40af-8c5c-488fc5c9a394', 'User', '7 km, trip to a supermaket', '1970-01-01T00:00:00Z'),
  ('9c04081e-4826-489b-b7b3-3734e4e49d67', '3ea8ea06-d2d2-40af-8c5c-488fc5c9a394', 'User', '0 km, I was lazy af', '1970-01-01T00:00:00Z'),
  ('fbc762dc-2895-45bc-ae56-a73a7fa6e8b5', '3a0c9a3b-bbec-4047-9822-1c4806c2a258', 'User', 'Checkout Coursera', '1970-01-01T00:00:00Z'),
  ('76ba6085-8d85-4c99-a6da-995f1ff4895b', '422eafc4-e488-47b3-9eb4-efd99162bd3d', 'User', 'See the movie', '1970-01-01T00:00:00Z'),
  ('1acf639d-ac92-499c-91af-8436fc431c10', '0f10f18a-bd51-4822-a44f-f5786baf5d07', 'User', 'Buy mint meat', '1970-01-01T00:00:00Z')
;

COMMIT;
