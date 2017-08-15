var searchdata = Object
function SForm(e) {
	
	if (e.preventDefault) e.preventDefault();
	
	var lat = document.getElementById("search").elements[0].value;
	var lng = document.getElementById("search").elements[1].value;
	var keyword = document.getElementById("search").elements[2].value;
	$.get("` + ngrokFrontend + `/search",{LAT:lat,LNG:lng,KEYWORD: keyword }, function(data){
        searchdata = data
        values=JSON.stringify(data);
	alert("Results:"+"\n"+values);
	  });

  return false;
};

function RForm(e) {
    if (e.preventDefault) e.preventDefault();
    var Parems =[query1,query2,query3]
	 query1 = document.getElementById("analysis").elements[0].value;
	 query2 = document.getElementById("analysis").elements[1].value;
     query3 = document.getElementById("analysis").elements[2].value;

     searchvalues=JSON.stringify(searchdata);
     message =[Parems,searchvalues]
     jsonmessage=JSON.stringify(message)
     $.ajax({
        url : "` + ngrokFrontend + `/analysis",
        type: "POST",
        data: jsonmessage,
        contentType: "application/json",
        dataType   : "json",
        success    : function(resultdata){
            	values=JSON.stringify(resultdata);
            alert("Results:"+"\n"+values);
        }
    });
  return false;
};


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
