-- =================================================================
-- FLASHQUEST - SCRIPT DE TESTES PARA POPULAR BANCO DE DADOS LOCAL
-- =================================================================

-- Limpar dados anteriores (CUIDADO!)
-- TRUNCATE TABLE user_response, comment_relationship, comment, 
--              question_set_question, question_set, 
--              question_topic, topic, 
--              question_source, answer, 
--              source_exam_instance, question, source, 
--              users, subject 
-- RESTART IDENTITY CASCADE;

-- =================================================================
-- 1. INSERIR USUÁRIOS
-- =================================================================
INSERT INTO users (id, name, email, password_hash) VALUES
(1, 'João Silva', 'joao.silva@test.com', 'hash_senha_123'),
(2, 'Maria Santos', 'maria.santos@test.com', 'hash_senha_456'),
(3, 'Pedro Costa', 'pedro.costa@test.com', 'hash_senha_789'),
(4, 'Ana Oliveira', 'ana.oliveira@test.com', 'hash_senha_000'),
(5, 'FlashQuest AI', 'ia@flashquest.com', 'hash_ia_gemini')
ON CONFLICT DO NOTHING;

-- =================================================================
-- 2. INSERIR DISCIPLINAS (SUBJECTS)
-- =================================================================
INSERT INTO subject (id, name) VALUES
(1, 'Direito Constitucional'),
(2, 'Direito Penal'),
(3, 'Direito Civil'),
(4, 'Direito Administrativo'),
(5, 'Direito do Trabalho'),
(6, 'Direito Tributário'),
(7, 'Geral')
ON CONFLICT DO NOTHING;

-- =================================================================
-- 3. INSERIR TÓPICOS (TOPICS)
-- =================================================================
INSERT INTO topic (id, subject_id, name) VALUES
(1, 1, 'Fundamentos da CF/88'),
(2, 1, 'Direitos e Garantias Fundamentais'),
(3, 2, 'Parte Geral do Código Penal'),
(4, 2, 'Crimes contra a Pessoa'),
(5, 3, 'Das Obrigações'),
(6, 3, 'Das Coisas'),
(7, 4, 'Princípios da Administração Pública'),
(8, 4, 'Atos Administrativos'),
(9, 5, 'Contrato de Trabalho'),
(10, 5, 'Direitos do Trabalhador')
ON CONFLICT DO NOTHING;

-- =================================================================
-- 4. INSERIR FONTES DE EXAMES
-- =================================================================
INSERT INTO source (id, name, type, created_at) VALUES
(1, 'OAB - Exame de Ordem Unificado', 'Exame', CURRENT_TIMESTAMP),
(2, 'Concursos Públicos - Magistratura', 'Exame', CURRENT_TIMESTAMP),
(3, 'CESPE/CEBRASPE', 'Banca', CURRENT_TIMESTAMP),
(4, 'FGV Projetos', 'Banca', CURRENT_TIMESTAMP),
(5, 'Estudo Livre', 'Customizado', CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

-- =================================================================
-- 5. INSERIR INSTÂNCIAS DE EXAME (EXAM INSTANCES)
-- =================================================================
INSERT INTO source_exam_instance (id, source_id, edition, phase, year, created_at) VALUES
(1, 1, 35, 1, 2024, CURRENT_TIMESTAMP),
(2, 1, 35, 2, 2024, CURRENT_TIMESTAMP),
(3, 2, 1, 1, 2024, CURRENT_TIMESTAMP),
(4, 3, 1, 1, 2023, CURRENT_TIMESTAMP),
(5, 5, 1, 1, 2025, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

-- =================================================================
-- 6. INSERIR QUESTÕES
-- =================================================================
INSERT INTO question (id, statement, subject_id, user_id, source_exam_instance_id, created_at, updated_at) VALUES
(1, 'A Constituição Federal de 1988 é considerada uma constituição:', 1, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 'Qual é o órgão máximo do judiciário no Brasil?', 1, 1, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, 'O homicídio é punido com pena de:', 2, 2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(4, 'Qual a diferença fundamental entre furto e roubo?', 2, 2, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(5, 'A propriedade é um direito garantido pela Constituição?', 3, 1, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(6, 'Qual desses princípios NÃO é princípio da Administração Pública?', 4, 3, 4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(7, 'O que é um ato administrativo válido?', 4, 3, 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(8, 'A CLT é o principal estatuto regulador das relações trabalhistas?', 5, 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(9, 'Qual é a jornada máxima de trabalho permitida por lei?', 5, 2, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(10, 'O ICMS é um tributo federal?', 6, 4, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

-- =================================================================
-- 7. INSERIR RESPOSTAS (ANSWERS)
-- =================================================================
INSERT INTO answer (id, id_question, text, is_correct) VALUES
-- Questão 1
(1, 1, 'Rígida', TRUE),
(2, 1, 'Flexível', FALSE),
(3, 1, 'Semi-rígida', FALSE),
(4, 1, 'Mutável', FALSE),

-- Questão 2
(5, 2, 'Supremo Tribunal Federal (STF)', TRUE),
(6, 2, 'Superior Tribunal de Justiça (STJ)', FALSE),
(7, 2, 'Tribunal de Justiça local', FALSE),
(8, 2, 'Conselho Nacional de Justiça (CNJ)', FALSE),

-- Questão 3
(9, 3, 'Reclusão de 6 a 20 anos', TRUE),
(10, 3, 'Detenção de 1 a 3 anos', FALSE),
(11, 3, 'Multa de até 50 salários mínimos', FALSE),
(12, 3, 'Reclusão de 12 a 30 anos', FALSE),

-- Questão 4
(13, 4, 'O roubo tem violência ou grave ameaça, o furto não', TRUE),
(14, 4, 'O furto é crime sem pena de prisão', FALSE),
(15, 4, 'O roubo é praticado apenas noite', FALSE),
(16, 4, 'Não há diferença legal entre eles', FALSE),

-- Questão 5
(17, 5, 'Sim, conforme art. 5º da CF/88', TRUE),
(18, 5, 'Não, é uma permissão do Estado', FALSE),
(19, 5, 'Apenas propriedade privada', FALSE),
(20, 5, 'Apenas propriedade pública', FALSE),

-- Questão 6
(21, 6, 'Legalidade', FALSE),
(22, 6, 'Moralidade', FALSE),
(23, 6, 'Publicidade', FALSE),
(24, 6, 'Lucro máximo', TRUE),

-- Questão 7
(25, 7, 'Aquele que respeita os pressupostos legais', TRUE),
(26, 7, 'Qualquer decisão de um servidor público', FALSE),
(27, 7, 'Um ato que beneficia a administração', FALSE),
(28, 7, 'Uma ordem do Presidente', FALSE),

-- Questão 8
(29, 8, 'Sim, é o principal código do direito do trabalho', TRUE),
(30, 8, 'Não, há apenas leis esparsas', FALSE),
(31, 8, 'Parcialmente, é complementada por decretos', FALSE),
(32, 8, 'Não mais, foi revogada em 2023', FALSE),

-- Questão 9
(33, 9, '8 horas diárias ou 44 horas semanais', TRUE),
(34, 9, '6 horas diárias', FALSE),
(35, 9, '10 horas diárias', FALSE),
(36, 9, 'Sem limite legal', FALSE),

-- Questão 10
(37, 10, 'Não, é um tributo estadual/municipal', TRUE),
(38, 10, 'Sim, gerenciado pela Receita Federal', FALSE),
(39, 10, 'Sim, é a principal arrecadação federal', FALSE),
(40, 10, 'É municipal apenas', FALSE)
ON CONFLICT DO NOTHING;

-- =================================================================
-- 8. INSERIR LISTAS DE EXERCÍCIOS (QUESTION SETS)
-- =================================================================
INSERT INTO question_set (id, user_id, type, name, description, is_private, created_at) VALUES
(1, 1, 'Simulado', 'Simulado OAB - Primeira Fase 2024', 'Questões de todas as disciplinas para primeira fase da OAB', FALSE, CURRENT_TIMESTAMP),
(2, 2, 'Lista de Exercícios', 'Direito Constitucional - Fundamentals', 'Exercícios básicos de Direito Constitucional', TRUE, CURRENT_TIMESTAMP),
(3, 1, 'Revisão', 'Revisão Direito Penal 2024', 'Questões para revisão geral de Direito Penal', FALSE, CURRENT_TIMESTAMP),
(4, 3, 'Prova', 'Prova Simulada Magistratura', 'Para concurso de magistratura estadual', TRUE, CURRENT_TIMESTAMP),
(5, 4, 'Estudo', 'Tópicos com dúvidas - Direito Civil', 'Questões em tópicos onde tenho dificuldade', TRUE, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

-- =================================================================
-- 9. ASSOCIAR QUESTÕES COM LISTAS (QUESTION SET QUESTIONS)
-- =================================================================
INSERT INTO question_set_question (question_set_id, question_id, position) VALUES
-- Simulado OAB
(1, 1, 1),
(1, 2, 2),
(1, 3, 3),
(1, 4, 4),
(1, 5, 5),
(1, 6, 6),
(1, 7, 7),
(1, 8, 8),
(1, 9, 9),
(1, 10, 10),

-- Lista Direito Constitucional
(2, 1, 1),
(2, 2, 2),

-- Revisão Direito Penal
(3, 3, 1),
(3, 4, 2),

-- Prova Simulada
(4, 1, 1),
(4, 5, 2),
(4, 8, 3),

-- Estudo Direito Civil
(5, 5, 1)
ON CONFLICT DO NOTHING;

-- =================================================================
-- 10. ASSOCIAR QUESTÕES COM TÓPICOS (QUESTION TOPICS)
-- =================================================================
INSERT INTO question_topic (question_id, topic_id) VALUES
(1, 1),
(2, 2),
(3, 3),
(4, 4),
(5, 6),
(6, 7),
(7, 8),
(8, 9),
(9, 10),
(10, 10)
ON CONFLICT DO NOTHING;

-- =================================================================
-- 11. ATUALIZAR SEQUÊNCIAS (para geração automática de IDs)
-- =================================================================
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
SELECT setval('subject_id_seq', (SELECT MAX(id) FROM subject));
SELECT setval('source_id_seq', (SELECT MAX(id) FROM source));
SELECT setval('source_exam_instance_id_seq', (SELECT MAX(id) FROM source_exam_instance));
SELECT setval('question_id_seq', (SELECT MAX(id) FROM question));
SELECT setval('answer_id_seq', (SELECT MAX(id) FROM answer));
SELECT setval('question_set_id_seq', (SELECT MAX(id) FROM question_set));

-- =================================================================
-- FIM DO SCRIPT
-- =================================================================
