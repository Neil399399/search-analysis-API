package main

var WebView = `
<!DOCTYPE html>
<html>
<body>

<script>

</script>

<form action="/search-analysis" method="get">
<h2>Location?<h2>
<h3>
LAT:<br>
  <input type="text" name="LAT" value="lat">
  <br>
LNG:<br>
  <input type="text" name="LNG" value="lng">
  <br><h3>
<h2>Want find?<h2>
<h3>
Keyword:<br>
  <input type="text" name="KEYWORD" value="word">
  <br><h3>
<h2>Query?<h2>
<h3>
  Query1:<br>
  <input type="text" name="analysis_word1" value="Comment">
  <br>
  Query2:<br>
  <input type="text" name="analysis_word2" value="Comment">
  <br>
   Query3:<br>
  <input type="text" name="analysis_word3" value="Comment">
  <br><h3>
  <br>
  <input type="submit" value="Submit">
</form> 
<p>If you click the "Submit" button, the form-data will be sent to a page called "/search-analysis".</p>

</body>
</html>
`
