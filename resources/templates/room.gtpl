{{define "title" }}
Room page
{{end}}

{{define "head" }}
{{end}}
{{define "content" }}
<h1>Chat</h1>
<div id='canvas'>
</div>
<div class="signout-div">
  <input type="text" id="message" autocomplete="off">
  <input class="btn btn-primary" type="submit" id='submit' value="Send" onclick="submitButton()">
</div>
{{end}}

{{define "footer" }}
{{end}}
