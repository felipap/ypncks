
ENDPOINT = "https://s3.amazonaws.com/ypncks/data.json"

$(function() {
    var $tableWrapper = $("#js-menu-table");
    
    function tablefy(data) {
        var $table = $("<table />");
        
        var colAttrs = [
            //{ id: "Id", name: "Id" },
            { id: "Name", name: "Name" },
            //{ id: "IdLocation", name: "IdLocation" },
            //{ id: "LocationCode", name: "LocationCode" },
            { id: "Location", name: "Location" },
            { id: "MealName", name: "MealName" },
            //{ id: "MealCode", name: "MealCode" },
            { id: "MenuDate", name: "MenuDate" },
            //{ id: "Course", name: "Course" },
            //{ id: "CourseCode", name: "CourseCode" },
            //{ id: "MenuItemId", name: "MenuItemId" },
            { id: "IsPar", name: "IsPar" },
            { id: "MealOpens", name: "MealOpens" },
            { id: "MealCloses", name: "MealCloses" },
            //{ id: "IsDefaultMeal", name: "IsDefaultMeal" },
            { id: "IsMenu", name: "IsMenu" },
        ];


        var $header = $("<tr>");
        colAttrs.forEach(function(col){
            $header.append($("<th>").html(col.name));
        });
        $table.append($header);

        data.forEach(function(rdata){
            var $tr = $("<tr>");
            colAttrs.forEach(function(col){
                var $td = $("<td>").html(rdata[col.id]);
                $tr.append($td);
            });
            $table.append($tr);
        })

        return $table;
    }

    $.getJSON(ENDPOINT, function (data) {
        var $table = tablefy(data);
        console.log($table);
        $tableWrapper.append($table);
    });
});
