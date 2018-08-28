-- phpMyAdmin SQL Dump
-- version 4.6.6deb4
-- https://www.phpmyadmin.net/
--
-- Хост: localhost:3306
-- Время создания: Авг 28 2018 г., 21:58
-- Версия сервера: 5.7.23
-- Версия PHP: 7.0.30-0+deb9u1

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- База данных: `CypherDesk-main`
--

-- --------------------------------------------------------

--
-- Структура таблицы `departments`
--

CREATE TABLE `departments` (
  `id` int(8) NOT NULL,
  `name` varchar(32) CHARACTER SET utf8 NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Дамп данных таблицы `departments`
--

INSERT INTO `departments` (`id`, `name`) VALUES
(1, 'admin'),
(2, 'kekLOL'),
(3, 'Лох');

-- --------------------------------------------------------

--
-- Структура таблицы `logs`
--

CREATE TABLE `logs` (
  `id` int(16) NOT NULL,
  `ticket` int(16) NOT NULL,
  `userFrom` int(16) NOT NULL,
  `userTo` int(16) NOT NULL,
  `action` varchar(16) CHARACTER SET utf8 NOT NULL,
  `time` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Структура таблицы `tickets`
--

CREATE TABLE `tickets` (
  `id` int(16) NOT NULL,
  `caption` varchar(32) CHARACTER SET utf8 NOT NULL,
  `description` text CHARACTER SET utf8 NOT NULL,
  `sender` int(16) NOT NULL,
  `status` varchar(16) CHARACTER SET utf8 NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Структура таблицы `users`
--

CREATE TABLE `users` (
  `id` int(16) NOT NULL,
  `login` varchar(32) CHARACTER SET utf8 NOT NULL,
  `pass` varchar(64) CHARACTER SET utf8 NOT NULL,
  `mail` varchar(64) CHARACTER SET utf8 DEFAULT NULL,
  `name` varchar(32) CHARACTER SET utf8 DEFAULT NULL,
  `surname` varchar(32) CHARACTER SET utf8 DEFAULT NULL,
  `partonymic` varchar(32) CHARACTER SET utf8 DEFAULT NULL,
  `recourse` varchar(16) CHARACTER SET utf8 DEFAULT NULL,
  `role` varchar(16) CHARACTER SET utf8 NOT NULL,
  `department` int(16) NOT NULL,
  `status` varchar(16) CHARACTER SET utf8 DEFAULT NULL,
  `activationKey` varchar(64) CHARACTER SET utf8 DEFAULT NULL,
  `activationType` int(3) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Дамп данных таблицы `users`
--

INSERT INTO `users` (`id`, `login`, `pass`, `mail`, `name`, `surname`, `partonymic`, `recourse`, `role`, `department`, `status`, `activationKey`, `activationType`) VALUES
(1, 'suriknik', '1860cf01768b32d9d32512a25c2789db88166cf54ec8cdac85df3b3a66f53a96', 'suriknik@bk.ru', 'Никита', 'Сурначев', 'Владимирович', 'док', 'admin', 1, NULL, '', NULL),
(11, 'surikon', '1860cf01768b32d9d32512a25c2789db88166cf54ec8cdac85df3b3a66f53a96', 'nikita.surnachev03@gmail.com', 'Никита', 'Сурначев', 'Владимирович', 'док', 'user', 2, '', '', NULL);

--
-- Индексы сохранённых таблиц
--

--
-- Индексы таблицы `departments`
--
ALTER TABLE `departments`
  ADD PRIMARY KEY (`id`);

--
-- Индексы таблицы `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT для сохранённых таблиц
--

--
-- AUTO_INCREMENT для таблицы `departments`
--
ALTER TABLE `departments`
  MODIFY `id` int(8) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;
--
-- AUTO_INCREMENT для таблицы `users`
--
ALTER TABLE `users`
  MODIFY `id` int(16) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=12;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
