package main

var WebView = `
<!DOCTYPE html>
<html>
<body>
<h1>Search and Analysis</h1>

<form id="search">
<h3>
LAT&nbsp /&nbsp LNG<br>
  <input type="text" name="LAT" value="lat">  <input type="text" name="LNG" value="lng">
  <br>
Keyword:<br>
  <input type="text" name="KEYWORD" value="word">  <button type="submit">Search!!</button> 
  <br><h3>

<hr color=#ff6600>
</form> 

<form id="analysis">
<h1>Analysis</h1>
<textarea rows="4" cols="50"></textarea>
<h3>
  Query:<br>
  <input type="text" name="analysis_word1" value="Comment1">  <input type="text" name="analysis_word2" value="Comment2"> <input type="text" name="analysis_word3" value="Comment3">
  <br>

 <button type="submit" value="Analysis">Analysis!!</button>
<hr color=#ff6600>
</form> 


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
 <button type="submit">Search and Analysis!!</button>
</form> 



<script
  src="https://code.jquery.com/jquery-3.2.1.min.js"
  integrity="sha256-hwg4gsxgFZhOsEEamdOYGBf13FyQuiTwlAQgxVSNgt4="
  crossorigin="anonymous"></script>

<script>
function SForm(e) {
	
	if (e.preventDefault) e.preventDefault();
	
	var lat = document.getElementById("search").elements[0].value;
	var lng = document.getElementById("search").elements[1].value;
	var keyword = document.getElementById("search").elements[2].value;
	$.get("` + ngrokFrontend + `/search",{LAT:lat,LNG:lng,KEYWORD: keyword }, function(data){
		values=JSON.stringify(data);
	alert("Results:"+"\n"+values);
		

	  });

  return false;
};


var xmlhttp = new XMLHttpRequest();












function SAForm(e) {
	if (e.preventDefault) e.preventDefault();
	
	var lat = document.getElementById("search-analysis").elements[0].value;
	var lng = document.getElementById("search-analysis").elements[1].value;
	var keyword = document.getElementById("search-analysis").elements[2].value;
	var query1 = document.getElementById("search-analysis").elements[3].value;
	var query2 = document.getElementById("search-analysis").elements[4].value;
	var query3 = document.getElementById("search-analysis").elements[5].value;
	
	$.get("` + ngrokFrontend + `/search-analysis",{LAT:lat,LNG:lng,KEYWORD: keyword ,analysis_word1:query1 ,analysis_word2: query2,analysis_word3: query3}, function(data){
	alert("Recomend:"+"\n"+data);
	  });
  return false;
};
var form = document.getElementById("search-analysis");
if (form.attachEvent) {
    form.attachEvent("submit", SAForm);
} else {
    form.addEventListener("submit", SAForm);
}

var form = document.getElementById("search");
if (form.attachEvent) {
    form.attachEvent("submit",  SForm);
} else  {
    form.addEventListener("submit",  SForm);
}
</script>


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
