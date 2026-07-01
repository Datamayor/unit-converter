package main

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

// CONVERSION LOGIC

func convertLength(value float64, from, to string) float64 {
	toMeters := map[string]float64{
		"millimeter": 0.001,
		"centimeter": 0.01,
		"meter":      1,
		"kilometer":  1000,
		"inch":       0.0254,
		"foot":       0.3048,
		"yard":       0.9144,
		"mile":       1609.344,
	}
	return value * toMeters[from] / toMeters[to]
}

func convertWeight(value float64, from, to string) float64 {
	toGrams := map[string]float64{
		"milligram": 0.001,
		"gram":      1,
		"kilogram":  1000,
		"ounce":     28.3495,
		"pound":     453.592,
	}
	return value * toGrams[from] / toGrams[to]
}

func convertTemperature(value float64, from, to string) float64 {
	// Step 1: to Celsius
	var celsius float64
	switch from {
	case "celsius":
		celsius = value
	case "fahrenheit":
		celsius = (value - 32) * 5 / 9
	case "kelvin":
		celsius = value - 273.15
	}
	// Step 2: from Celsius to target
	switch to {
	case "celsius":
		return celsius
	case "fahrenheit":
		return celsius*9/5 + 32
	case "kelvin":
		return celsius + 273.15
	}
	return 0
}

func round(value float64) string {
	rounded := math.Round(value*100000) / 100000
	return strconv.FormatFloat(rounded, 'f', -1, 64)
}

// TEMPLATE DATA

type PageData struct {
	Page      string
	Body      template.HTML
	Result    string
	FromValue string
	FromUnit  string
	ToUnit    string
	Error     string
}

// TEMPLATES

var layoutTmpl = template.Must(template.New("layout").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Unit Converter</title>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body { font-family: system-ui, sans-serif; background: #f5f5f5; color: #111; min-height: 100vh; }
    nav { background: #fff; border-bottom: 1px solid #e5e5e5; padding: 16px 32px; display: flex; gap: 8px; align-items: center; }
    nav .logo { font-weight: 600; font-size: 16px; color: #111; margin-right: 16px; }
    nav a { text-decoration: none; color: #555; font-size: 14px; padding: 7px 16px; border-radius: 6px; }
    nav a:hover { background: #f0f0f0; color: #111; }
    nav a.active { background: #111; color: #fff; }
    .container { max-width: 500px; margin: 48px auto; padding: 0 16px; }
    h1 { font-size: 22px; font-weight: 600; margin-bottom: 24px; }
    .card { background: #fff; border: 1px solid #e5e5e5; border-radius: 12px; padding: 28px; }
    label { font-size: 13px; color: #555; display: block; margin-bottom: 6px; }
    input, select { width: 100%; padding: 10px 12px; border: 1px solid #d5d5d5; border-radius: 8px; font-size: 15px; margin-bottom: 16px; font-family: inherit; }
    input:focus, select:focus { outline: none; border-color: #111; }
    button { width: 100%; padding: 11px; background: #111; color: #fff; border: none; border-radius: 8px; font-size: 15px; cursor: pointer; margin-top: 4px; font-family: inherit; }
    button:hover { background: #333; }
    .row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
    .result { margin-top: 20px; padding: 20px; background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 8px; text-align: center; }
    .result .value { font-size: 28px; font-weight: 600; color: #15803d; }
    .result .label { font-size: 13px; color: #555; margin-top: 6px; }
    .error { margin-top: 16px; padding: 12px; background: #fef2f2; border: 1px solid #fecaca; border-radius: 8px; color: #dc2626; font-size: 14px; }
  </style>
</head>
<body>
  <nav>
    <span class="logo">Unit Converter</span>
    <a href="/length"      {{if eq .Page "length"}}class="active"{{end}}>Length</a>
    <a href="/weight"      {{if eq .Page "weight"}}class="active"{{end}}>Weight</a>
    <a href="/temperature" {{if eq .Page "temperature"}}class="active"{{end}}>Temperature</a>
  </nav>
  <div class="container">
    {{.Body}}
  </div>
</body>
</html>`))

var lengthTmpl = template.Must(template.New("length").Parse(`
<h1>Length Converter</h1>
<div class="card">
  <form method="POST" action="/length">
    <label>Value</label>
    <input type="number" name="value" step="any" placeholder="Enter a value" value="{{.FromValue}}" required>
    <div class="row">
      <div>
        <label>From</label>
        <select name="from">
          <option value="millimeter" {{if eq .FromUnit "millimeter"}}selected{{end}}>Millimeter</option>
          <option value="centimeter" {{if eq .FromUnit "centimeter"}}selected{{end}}>Centimeter</option>
          <option value="meter"      {{if eq .FromUnit "meter"}}selected{{end}}>Meter</option>
          <option value="kilometer"  {{if eq .FromUnit "kilometer"}}selected{{end}}>Kilometer</option>
          <option value="inch"       {{if eq .FromUnit "inch"}}selected{{end}}>Inch</option>
          <option value="foot"       {{if eq .FromUnit "foot"}}selected{{end}}>Foot</option>
          <option value="yard"       {{if eq .FromUnit "yard"}}selected{{end}}>Yard</option>
          <option value="mile"       {{if eq .FromUnit "mile"}}selected{{end}}>Mile</option>
        </select>
      </div>
      <div>
        <label>To</label>
        <select name="to">
          <option value="millimeter" {{if eq .ToUnit "millimeter"}}selected{{end}}>Millimeter</option>
          <option value="centimeter" {{if eq .ToUnit "centimeter"}}selected{{end}}>Centimeter</option>
          <option value="meter"      {{if eq .ToUnit "meter"}}selected{{end}}>Meter</option>
          <option value="kilometer"  {{if eq .ToUnit "kilometer"}}selected{{end}}>Kilometer</option>
          <option value="inch"       {{if eq .ToUnit "inch"}}selected{{end}}>Inch</option>
          <option value="foot"       {{if eq .ToUnit "foot"}}selected{{end}}>Foot</option>
          <option value="yard"       {{if eq .ToUnit "yard"}}selected{{end}}>Yard</option>
          <option value="mile"       {{if eq .ToUnit "mile"}}selected{{end}}>Mile</option>
        </select>
      </div>
    </div>
    <button type="submit">Convert</button>
  </form>
  {{if .Result}}
  <div class="result">
    <div class="value">{{.Result}} {{.ToUnit}}</div>
    <div class="label">{{.FromValue}} {{.FromUnit}} = {{.Result}} {{.ToUnit}}</div>
  </div>
  {{end}}
  {{if .Error}}<div class="error">{{.Error}}</div>{{end}}
</div>`))

var weightTmpl = template.Must(template.New("weight").Parse(`
<h1>Weight Converter</h1>
<div class="card">
  <form method="POST" action="/weight">
    <label>Value</label>
    <input type="number" name="value" step="any" placeholder="Enter a value" value="{{.FromValue}}" required>
    <div class="row">
      <div>
        <label>From</label>
        <select name="from">
          <option value="milligram" {{if eq .FromUnit "milligram"}}selected{{end}}>Milligram</option>
          <option value="gram"      {{if eq .FromUnit "gram"}}selected{{end}}>Gram</option>
          <option value="kilogram"  {{if eq .FromUnit "kilogram"}}selected{{end}}>Kilogram</option>
          <option value="ounce"     {{if eq .FromUnit "ounce"}}selected{{end}}>Ounce</option>
          <option value="pound"     {{if eq .FromUnit "pound"}}selected{{end}}>Pound</option>
        </select>
      </div>
      <div>
        <label>To</label>
        <select name="to">
          <option value="milligram" {{if eq .ToUnit "milligram"}}selected{{end}}>Milligram</option>
          <option value="gram"      {{if eq .ToUnit "gram"}}selected{{end}}>Gram</option>
          <option value="kilogram"  {{if eq .ToUnit "kilogram"}}selected{{end}}>Kilogram</option>
          <option value="ounce"     {{if eq .ToUnit "ounce"}}selected{{end}}>Ounce</option>
          <option value="pound"     {{if eq .ToUnit "pound"}}selected{{end}}>Pound</option>
        </select>
      </div>
    </div>
    <button type="submit">Convert</button>
  </form>
  {{if .Result}}
  <div class="result">
    <div class="value">{{.Result}} {{.ToUnit}}</div>
    <div class="label">{{.FromValue}} {{.FromUnit}} = {{.Result}} {{.ToUnit}}</div>
  </div>
  {{end}}
  {{if .Error}}<div class="error">{{.Error}}</div>{{end}}
</div>`))

var temperatureTmpl = template.Must(template.New("temperature").Parse(`
<h1>Temperature Converter</h1>
<div class="card">
  <form method="POST" action="/temperature">
    <label>Value</label>
    <input type="number" name="value" step="any" placeholder="Enter a value" value="{{.FromValue}}" required>
    <div class="row">
      <div>
        <label>From</label>
        <select name="from">
          <option value="celsius"    {{if eq .FromUnit "celsius"}}selected{{end}}>Celsius</option>
          <option value="fahrenheit" {{if eq .FromUnit "fahrenheit"}}selected{{end}}>Fahrenheit</option>
          <option value="kelvin"     {{if eq .FromUnit "kelvin"}}selected{{end}}>Kelvin</option>
        </select>
      </div>
      <div>
        <label>To</label>
        <select name="to">
          <option value="celsius"    {{if eq .ToUnit "celsius"}}selected{{end}}>Celsius</option>
          <option value="fahrenheit" {{if eq .ToUnit "fahrenheit"}}selected{{end}}>Fahrenheit</option>
          <option value="kelvin"     {{if eq .ToUnit "kelvin"}}selected{{end}}>Kelvin</option>
        </select>
      </div>
    </div>
    <button type="submit">Convert</button>
  </form>
  {{if .Result}}
  <div class="result">
    <div class="value">{{.Result}} {{.ToUnit}}</div>
    <div class="label">{{.FromValue}} {{.FromUnit}} = {{.Result}} {{.ToUnit}}</div>
  </div>
  {{end}}
  {{if .Error}}<div class="error">{{.Error}}</div>{{end}}
</div>`))

// RENDER HELPER

func render(w http.ResponseWriter, page string, bodyTmpl *template.Template, data PageData) {
	// Render the inner body first into a buffer
	var buf bytes.Buffer
	bodyTmpl.Execute(&buf, data)

	// Now render the layout with the body inside
	data.Page = page
	data.Body = template.HTML(buf.String())
	layoutTmpl.Execute(w, data)
}

// HANDLERS

func handleLength(w http.ResponseWriter, r *http.Request) {
	data := PageData{FromUnit: "meter", ToUnit: "kilometer"}

	if r.Method == http.MethodPost {
		r.ParseForm()
		valueStr := r.FormValue("value")
		from := r.FormValue("from")
		to := r.FormValue("to")
		data.FromValue = valueStr
		data.FromUnit = from
		data.ToUnit = to

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			data.Error = "Please enter a valid number."
		} else {
			data.Result = round(convertLength(value, from, to))
		}
	}

	render(w, "length", lengthTmpl, data)
}

func handleWeight(w http.ResponseWriter, r *http.Request) {
	data := PageData{FromUnit: "kilogram", ToUnit: "pound"}

	if r.Method == http.MethodPost {
		r.ParseForm()
		valueStr := r.FormValue("value")
		from := r.FormValue("from")
		to := r.FormValue("to")
		data.FromValue = valueStr
		data.FromUnit = from
		data.ToUnit = to

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			data.Error = "Please enter a valid number."
		} else {
			data.Result = round(convertWeight(value, from, to))
		}
	}

	render(w, "weight", weightTmpl, data)
}

func handleTemperature(w http.ResponseWriter, r *http.Request) {
	data := PageData{FromUnit: "celsius", ToUnit: "fahrenheit"}

	if r.Method == http.MethodPost {
		r.ParseForm()
		valueStr := r.FormValue("value")
		from := r.FormValue("from")
		to := r.FormValue("to")
		data.FromValue = valueStr
		data.FromUnit = from
		data.ToUnit = to

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			data.Error = "Please enter a valid number."
		} else {
			data.Result = round(convertTemperature(value, from, to))
		}
	}

	render(w, "temperature", temperatureTmpl, data)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/length", http.StatusFound)
	})
	http.HandleFunc("/length", handleLength)
	http.HandleFunc("/weight", handleWeight)
	http.HandleFunc("/temperature", handleTemperature)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
