import {reduxRender} from "../../test-utils.tsx";
import SongDrawer from "./SongDrawer.tsx";
import {screen} from "@testing-library/react";
import {http, HttpResponse} from "msw";
import Song from "../../types/models/Song.ts";
import {setupServer} from "msw/node";

describe('Song Drawer', () => {
  const song: Song = {
    id: '1',
    title: "Justice for all",
    description: '',
    isRecorded: false,
    sections: []
  }

  const handlers = [
    http.get('/songs/1', async () => {
      return HttpResponse.json(song)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and display the loader', () => {
    reduxRender(<SongDrawer opened={true} close={() => {}} />)

    expect(screen.getByTestId('song-drawer-loader')).toBeInTheDocument()
  })

  it('should display song details when the songId exists', async () => {
    reduxRender(<SongDrawer opened={true} close={() => {}} />, { songs: { songId: '1' } })

    expect(await screen.findByText(song.title)).toBeInTheDocument()

    expect(screen.queryByTestId('song-drawer-loader')).not.toBeInTheDocument()
  })
})
