# ЧМ-2026 · FIFA World Cup 2026 Match Tracker

A small, fast **macOS desktop app** for following the 2026 FIFA World Cup: browse the full schedule in your own timezone, enter scores, make predictions, and watch the knockout bracket and group tables fill in automatically — all offline, with results saved next to the app.

Built with [Wails](https://wails.io) (Go + vanilla JS, no frameworks). Dark UI. 6 languages.

> 🇬🇧 English below · 🇷🇺 Русская версия ниже

---

## 🇬🇧 English

### Features
- **Schedule** with a date carousel that auto-scrolls to today; live matches show an **Online** badge and a blinking score.
- **Score entry** with instant autosave (no buttons, no confirmations).
- **Predictions** — enter your forecast above each score; a stats page tracks exact hits / correct outcomes / misses and points.
- **Knockout bracket** — a real two-sided tree that fills in automatically from results, with official FIFA naming and hover tooltips for every slot.
- **Group tables** with the official **FIFA 2026 tie-breakers** (head-to-head → goal difference → goals scored …) and best-third-placed-teams ranking + allocation.
- **Timezone** selector — all kick-off times are recalculated for the zone you pick (DST-aware).
- **6 languages**: English, Deutsch, Français, Español, Русский, 繁體中文.
- Everything persists between launches (`results.json` + `settings.json` next to the binary).

### Requirements
- macOS (Apple Silicon / M-series)
- [Go](https://go.dev) ≥ 1.23
- [Node.js](https://nodejs.org)
- [Wails CLI v2](https://wails.io): `go install github.com/wailsapp/wails/v2/cmd/wails@latest` (then ensure `$(go env GOPATH)/bin` is on your `PATH`)

### Run from source
```bash
wails dev
```

### Build (Apple Silicon)
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
wails build -platform darwin/arm64
```
The unsigned binary is at `build/bin/wc2026.app`. The app is ad-hoc signed, so on first launch macOS may warn about an unidentified developer — **right-click → Open**, or run:
```bash
xattr -dr com.apple.quarantine build/bin/wc2026.app
```
> If `wails build` fails at the self-signing step with *“resource fork … not allowed”* (happens when the project lives in an iCloud-synced folder like `~/Documents`), strip the extended attributes and re-sign from a non-synced folder:
> ```bash
> ditto --norsrc --noextattr --noacl build/bin/wc2026.app /tmp/wc2026.app
> codesign --force --deep -s - /tmp/wc2026.app
> open /tmp/wc2026.app
> ```

### Data
Match results and your predictions live in `results.json`, settings (timezone + language) in `settings.json`, both written next to the executable. Back them up or sync them (e.g. via iCloud) as you like.

### Releases & distribution
Prebuilt apps are published on the **[Releases](https://github.com/kh-marik/wc2026/releases)** page; forkers build their own. Versions follow [SemVer](https://semver.org) with git tags `vX.Y.Z`. Maintainer steps per release:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
wails build -platform darwin/arm64
xattr -cr build/bin/wc2026.app && codesign --force --deep -s - build/bin/wc2026.app
ditto -c -k --keepParent build/bin/wc2026.app wc2026-vX.Y.Z-macos-arm64.zip
git tag vX.Y.Z && git push origin vX.Y.Z
gh release create vX.Y.Z wc2026-vX.Y.Z-macos-arm64.zip -t "vX.Y.Z" -n "Changelog…"
```
(The `.app` must be zipped with `ditto -c -k --keepParent` to preserve the bundle and signature. For a seamless install on other Macs, sign with an Apple Developer ID and notarize.)

### License
[MIT](LICENSE) — fork, modify and distribute freely.

---

## 🇷🇺 Русский

**ЧМ-2026** — небольшое десктоп-приложение для macOS, чтобы следить за чемпионатом мира 2026: всё расписание в вашем часовом поясе, ввод счёта, прогнозы, а сетка плей-офф и таблицы групп заполняются автоматически. Работает офлайн, результаты сохраняются рядом с приложением. Сделано на [Wails](https://wails.io) (Go + ванильный JS). Тёмная тема, 6 языков.

### Возможности
- **Расписание** с каруселью дат и автопрокруткой к сегодняшнему дню; у идущих матчей — бейдж **Онлайн** и моргающий счёт.
- **Ввод счёта** с мгновенным автосохранением (без кнопок и подтверждений).
- **Прогнозы** — вписываете прогноз над счётом; на странице статистики считаются точные попадания / угаданные исходы / промахи и очки.
- **Сетка плей-офф** — настоящее двустороннее дерево, заполняется автоматически по результатам; официальные названия FIFA и подсказки по каждому слоту.
- **Таблицы групп** по официальным **тай-брейкам FIFA 2026** (личные встречи → разница мячей → забитые …) с ранжированием и распределением лучших третьих мест.
- **Часовой пояс** — все времена начала пересчитываются под выбранную зону (с учётом летнего времени).
- **6 языков**: English, Deutsch, Français, Español, Русский, 繁體中文.
- Всё сохраняется между запусками (`results.json` + `settings.json` рядом с бинарником).

### Требования
- macOS (Apple Silicon)
- [Go](https://go.dev) ≥ 1.23, [Node.js](https://nodejs.org)
- [Wails CLI v2](https://wails.io): `go install github.com/wailsapp/wails/v2/cmd/wails@latest` (и добавьте `$(go env GOPATH)/bin` в `PATH`)

### Запуск из исходников
```bash
wails dev
```

### Сборка (Apple Silicon)
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
wails build -platform darwin/arm64
```
Приложение подписано ad-hoc, поэтому при первом запуске macOS может предупредить о «неустановленном разработчике» — **правый клик → «Открыть»**, либо `xattr -dr com.apple.quarantine build/bin/wc2026.app`.
> Если `wails build` падает на шаге подписи с *«resource fork … not allowed»* (бывает, когда проект лежит в iCloud-папке вроде `~/Documents`) — соберите, затем очистите атрибуты и переподпишите из несинхронизируемой папки:
> ```bash
> ditto --norsrc --noextattr --noacl build/bin/wc2026.app /tmp/wc2026.app
> codesign --force --deep -s - /tmp/wc2026.app && open /tmp/wc2026.app
> ```

### Данные
Результаты и прогнозы — в `results.json`, настройки (зона + язык) — в `settings.json`, оба рядом с исполняемым файлом.

### Релизы
Готовые сборки — на странице **[Releases](https://github.com/kh-marik/wc2026/releases)**. Версии — по [SemVer](https://semver.org), теги `vX.Y.Z`. Команды для релиза — см. раздел *Releases & distribution* выше.

### Лицензия
[MIT](LICENSE) — форкайте, меняйте и распространяйте свободно.
