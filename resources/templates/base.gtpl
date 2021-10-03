<!DOCTYPE html>
<html lang="en">
  <head>
    <title>{{block "title" .}}Simple chat{{end}}</title>
    <meta charset="UTF-8">
    <link href="/static/style.css" rel="stylesheet">
    <script src="/static/index.js"></script>
    {{block "head" .}}{{end}}
  </head>
  <body>
    Hello
    {{block "content" .}}{{end}}
    <footer>
      <span>Property of Vitalii Drevenchuk 2021</span>
      {{block "footer" .}}{{end}}
    </footer>
  </body>
</html>
