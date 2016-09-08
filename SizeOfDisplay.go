package main

var SizeOfDisplay = `{{ if .Error }}
    <div class="bs-callout bs-callout-danger">
        <h4>Parsing error</h4>
        <p>{{ .Error }}</p>
    </div>
{{ else }}
    {{ with .Result }}
          <h3>Type size: {{ .Sizeof }}</h3>
        {{ if .IsFixed }}
            <div class="bs-callout bs-callout-info">
                <h4>Explanation</h4>
                <p>Your type is {{ .Name }} and always has fixed sized, no matter how it was defined.</p>
            </div>
        {{ end }}
        {{ if .IsArray }}
            <div class="bs-callout bs-callout-info">
                <h4>Explanation</h4>
                <p>Your type is array and its size is a product of its length with its underlying type size.</p>
                <p>Formula: <code>Sizeof([N]Type) = N * Sizeof(Type)</code></p>
            </div>
        {{ end }}
        {{ if .IsStruct }}
            <div class="bs-callout bs-callout-info">
                <h4>Explanation</h4>
                <p>
                   Your type is struct and its size depends on how underlying types were defined.
                   Two structs with same fields but different order may have different sizes.
                   Size of struct is counted accordingly with padding and alignment rules.
               </p>
           </div>
           {{ if .Details }}
                <h3>Struct alignment: {{ .Alignof }}</h3>
                <div class="table-responsive">
                    <table>
                        <tr>
                            <th>Fields</th>
                            <th>Aligment</th>
                        </tr>
                        {{ range $row := .Details }}
                            <tr>
                                <td>{{ $row.Name }}</td>
                                <td>
                                    {{ $len := $row.Chunks | len }}
                                    {{ range $i, $ch := $row.Chunks }}
                                        {{ if and (gt $len 4) (unvischunk $i $len) }}
                                        {{ else }}
                                            {{ if and (gt $len 4) (eq $i 2) }}...<br/>{{ else }}
                                                {{ range $filled := $ch.Cells }}
                                                    <div class="chnk{{ if $ch.IsPadding }} pad{{ end }}{{ if not $filled }} empty{{ end }}"></div>
                                                {{ end }}
                                                <br/>
                                            {{ end }}
                                        {{ end }}
                                    {{ end }}
                                </td>
                                <td>{{ if (gt $len 4) }}{{ $len }} total{{ end }}</td>
                            </tr>
                        {{ end }}
                    </table>
                </div>
            {{ end }}
        {{ end }}
    {{ end }}
{{ end }}
`
