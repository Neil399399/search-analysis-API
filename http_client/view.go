package main

var WebView = `
<!DOCTYPE html>
<html>
<body>
<h1>Search</h1>

<form id="search">
<h3>
LAT&nbsp /&nbsp LNG<br>
  <input type="text" name="LAT" value="25.03978">  <input type="text" name="LNG" value="121.548495">
  <br>
Keyword:<br>
  <input type="text" name="KEYWORD" value="海鮮餐廳">  <button type="submit">Search!!</button> 
  <br><h3>

<hr color=#ff6600>
</form> 

<h1>Analysis</h1>
<form id="analysis">
<h3>
  Query:<br>
  <input type="text" name="analysis_word1" value="新鮮">  <input type="text" name="analysis_word2" value="好吃"> <input type="text" name="analysis_word3" value="便宜">
  <br>
<textarea rows="4" cols="50"></textarea>
 <button type="submit" value="Analysis">Analysis!!</button>
<hr color=#ff6600>
</form> 

<h1>Search and Analysis</h1>
<form id="search-analysis">
<h3>
LAT&nbsp /&nbsp LNG<br>
  <input type="text" name="LAT" value="lat">  <input type="text" name="LNG" value="lng">
  <br>
Keyword:<br>
  <input type="text" name="KEYWORD" value="word">
  <br>
  Query:<br>
  <input type="text" name="analysis_word1" value="Comment1">  <input type="text" name="analysis_word2" value="Comment2"> <input type="text" name="analysis_word3" value="Comment3">
  <h3>
 <button type="submit">Search and Analysis!!</button>
</form> 



<script
  src="https://code.jquery.com/jquery-3.2.1.min.js"
  integrity="sha256-hwg4gsxgFZhOsEEamdOYGBf13FyQuiTwlAQgxVSNgt4="
  crossorigin="anonymous"></script>

<script  src="static/API.js" type="text/javascript"></script>


</body>
</html>`

/*
<form id="search-analysis">
<h2>Location<h2>
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
<h2>Query<h2>
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
 <button type="submit">Go</button>
</form>




*/
