setInterval(function(){
   callServer();
}, 3000);

var crossings = Array(5).fill(0).map(()=>Array(5).fill(0));
              // N-Exit, NE-Exit, SE-Exit, SV-Exit, NV-Exit
  // N-Entry      0       0         0       0         0
  // NE-Entry     0       0         0       0         0
  // SE-Entry     0       0         0       0         0
  // SV-Entry     0       0         0       0         0
  // NV-Entry     0       0         0       0         0

var lengthEvents = 0;

function getLCIdxFromString(crossing) {
  let pos = crossing.search("-");
  let prefix;
  
  if (pos == -1) {
    prefix = crossing.substring(0, crossing.length);
  } else {
    prefix = crossing.substring(0, pos);
  }

  if (0 == prefix.localeCompare("N")) {
    return 0;
  } else if (0 == prefix.localeCompare("NE")) {
    return 1;
  } else if (0 == prefix.localeCompare("SE")) {
    return 2;
  } else if (0 == prefix.localeCompare("SV")) {
    return 3;
  } else if (0 == prefix.localeCompare("NV")) {
    return 4;
  }
}

function callServer(){

  const myHeaders = new Headers();
  myHeaders.append('Content-Type', 'application/json');
  const req = new Request('http://192.168.50.10:8080/events', {
    method:'GET',
    headers: myHeaders,
    mode: 'no-cors',
    cache: 'default',
  });

  fetch(req)
  .then(response => response.json())
  .then(data => {
    //console.log(data)
    lengthEvents = 0
    let dataObj = JSON.parse(data);
    if (!(Symbol.iterator in Object(dataObj.Events))) {
      console.log("not iterable");
      return;
    }
    crossings = Array(5).fill(0).map(()=>Array(5).fill(0));
    for (const el of dataObj.Events) {
      // console.log(el.event);
      // console.log(getLCIdxFromString(el.event.entry)); console.log(getLCIdxFromString(el.event.exit));
      crossings[getLCIdxFromString(el.event.entry)][getLCIdxFromString(el.event.exit)]++;
      lengthEvents++;
    }
    setEventsLength();
    buildChordDiagram();
  })
  .catch(console.error);
}

function setEventsLength()
{
  d3.select("h1").text(`Kafka Messages: ${lengthEvents}`);
}

function buildChordDiagram()
{
  const width = 964;
  const height = 964; 
  const innerRadius = 352;
  const outerRadius = 362;
  formatValue = x => `${x.toFixed(0)} vehicles`

  d3.select("svg").remove()

  // create the svg area
  const svg = d3.select("#my_dataviz")
    .append("svg")
      .attr("viewBox", [-width / 2, -height / 2, width, height])
      .attr("font-size", 14)
      .attr("font-family", "sans-serif")
      .style("width", "100%")
      .style("height", "auto")
    .append("g")
      .attr("transform", "rotate(-70)")

  // create a matrix
  // const matrix = [
  //   [0,  5871, 8916, 2868, 19],
  //   [ 1951, 0, 2060, 6171, 20],
  //   [ 8010, 16145, 0, 8045, 21],
  //   [ 1013,   990,  940, 0, 22],
  //   [ 8010, 16145, 0, 8045, 23]
  // ];

  data = {
    matrix: crossings,
    indexByName: {"N":0, "NE":1, "SE":2, "SV":3, "NV":4},
    nameByIndex: {0:"N", 1:"NE", 2:"SE", 3:"SV", 4:"NV"},
  };

  // 4 groups, so create a vector of 4 colors
  const colors = [ "#335c67", "#fff3b0", "#e09f3e", "#9e2a2b", "#540b0e"]

  // give this matrix to d3.chord(): it will calculates all the info we need to draw arc and ribbon
  const res = d3.chord()
      .padAngle(0.05)
      .sortSubgroups(d3.descending)
      (data.matrix)

  // add the groups on the outer part of the circle
  const group = svg.datum(res)
    .append("g")
    .selectAll("g")
    .data(function(d) { return d.groups; })
    .join("g");

  group.append("path")
      .style("fill", (d,i) => colors[i])
      .style("stroke", (d,i) => colors[i])
      .attr("d", d3.arc()
        .innerRadius(innerRadius)
        .outerRadius(outerRadius)
      )
      .attr("index", (d,i) => i)
      .on("mouseover", onMouseOver)
      .on("mouseout", onMouseOut);

  // Add the links between groups
  const links = svg.datum(res)
    .append("g")
    .selectAll("path")
    .data(function(d) { return d; })
    .enter()
    .append("path")
      // .attr("d", d3.ribbon()
      //   .radius(innerRadius)
      // )
      .attr("d", d3.ribbonArrow()
        .radius(innerRadius - 0.5)
        .padAngle(1/innerRadius)
      )
      .style("fill", function(d){ return(colors[d.source.index]) })
      .style("stroke", function(d){ return(d3.rgb(colors[d.source.index]).darker()) })
      .attr("index", (d) => d.source.index)
      .on("mouseover", onMouseOver) 
      .on("mouseout", onMouseOut);

  links.append("title")
        .text(d => `${data.nameByIndex[d.source.index]} => ${data.nameByIndex[d.target.index]}: ${formatValue(d.source.value)}`);

  group.append("text")
    .each(d => { d.angle = (d.startAngle + d.endAngle) / 2; })
    .attr("dy", ".35em")
    .attr("transform", d => `
      rotate(${(d.angle * 180 / Math.PI - 90)})
      translate(${outerRadius + 26})
      ${d.angle > Math.PI ? "rotate(180)" : ""}
    `)
    .attr("text-anchor", d => d.angle > Math.PI ? "end" : null)
    .style("fill", "rgb(241, 241, 241)")
    .text(d => data.nameByIndex[d.index]);

  group.append("title")
  .text(d => `${data.nameByIndex[d.index]}
  Entry: ${formatValue(d3.sum(data.matrix[d.index]))}
  Exit:  ${formatValue(d3.sum(data.matrix, row => row[d.index]))}`);
  
  function onMouseOver(selected) {
    group      
      .filter( d => d.index !== parseInt(selected.target.attributes.index.value))
      .style("opacity", 0.3);
    
    links
      .filter( d => d.source.index !== parseInt(selected.target.attributes.index.value))
      .style("opacity", 0.3);
  }
    
  function onMouseOut() {
    group.style("opacity", 1);
    links.style("opacity", 1);
  }
}

callServer();