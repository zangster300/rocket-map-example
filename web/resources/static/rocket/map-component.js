import { rocket } from 'datastar-pro'

const { default: maplibregl } = await import(
    'https://cdn.jsdelivr.net/npm/maplibre-gl@5.18.0/+esm'
)


rocket("map-component", {
    mode: 'light',
    props: ({ json }) => ({
        markers: json.default(() => [])
    }),
    setup: ({ props, observeProps }) => {
        const DARK_STYLE_URL = "https://basemaps.cartocdn.com/gl/dark-matter-gl-style/style.json"
        const LIGHT_STYLE_URL = "https://basemaps.cartocdn.com/gl/positron-gl-style/style.json"


        const prefersColorSchemeDark = window.matchMedia("(prefers-color-scheme: dark)")
        let styleURL = prefersColorSchemeDark.matches ? DARK_STYLE_URL : LIGHT_STYLE_URL

        let map = new maplibregl.Map({
            container: 'map', // container id
            style: styleURL,
            center: [0, 0], // starting position [lng, lat]
            zoom: 0, // starting zoom
            renderWorldCopies: false,
            attributionControl: false
        });

        prefersColorSchemeDark.addEventListener('change', (e) => {
            if (e.matches) {
                map.setStyle(DARK_STYLE_URL)
            } else {
                map.setStyle(LIGHT_STYLE_URL)
            }
        })



        // helper to create map markers
        function createMarker(username, userId, avatarHash) {
            const markerContainer = document.createElement('div')
            markerContainer.classList.add('marker-container')

            const markerElement = document.createElement('div')
            markerElement.classList.add('marker')
            markerElement.style.backgroundImage = `url("https://cdn.discordapp.com/avatars/${userId}/${avatarHash}.png")`

            markerContainer.appendChild(markerElement)

            return markerContainer
        }

        // watch pattern
        // https://data-star.dev/examples/rocket_globe#6d88edf638107c9a_line_29
        let prevMarkersStr = ''
        let prevMarkers = []

        observeProps(() => {
            if (!Array.isArray(props.markers)) return
            const markersStr = JSON.stringify(props.markers || [])
            if (prevMarkers.length > 0) {
                for (let prevMarker of prevMarkers) {
                    prevMarker.remove()
                }
            }
            if (markersStr !== prevMarkersStr) {
                prevMarkersStr = markersStr
                const markers = JSON.parse(markersStr)

                for (let marker of markers) {
                    const mapMarker = new maplibregl.Marker({ element: createMarker(marker.user.username, marker.user.id, marker.user.avatar) })
                        .setLngLat([marker.coordinates[0], marker.coordinates[1]])
                        .addTo(map)
                    prevMarkers.push(mapMarker)
                }
            }
        })
    }
})