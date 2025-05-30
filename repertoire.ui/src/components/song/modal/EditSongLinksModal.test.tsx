import { emptySong, reduxRender, withToastify } from '../../../test-utils.tsx'
import Song from '../../../types/models/Song.ts'
import { setupServer } from 'msw/node'
import { act, screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { UpdateSongRequest } from '../../../types/requests/SongRequests.ts'
import EditSongLinksModal from './EditSongLinksModal.tsx'

describe('Edit Song Links Modal', () => {
  const song: Song = {
    ...emptySong,
    id: 'some-id',
    youtubeLink: 'https://www.youtube.com/watch?v=VUNPqDXxG3U',
    songsterrLink: 'https://www.songsterr.com/a/wsa/metallica-nothing-else-matters-tab-s439171'
  }

  const newYoutubeLink = 'https://www.youtube.com/watch?v=2z55Cx5fcO4'
  const newSongsterrLink = 'https://www.songsterr.com/a/wsa/metallica-master-of-puppets-tab-s455118'

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<EditSongLinksModal opened={true} onClose={() => {}} song={song} />)

    expect(screen.getByRole('dialog', { name: /edit song links/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /edit song links/i })).toBeInTheDocument()

    expect(screen.getByRole('textbox', { name: /youtube/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /youtube/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /youtube/i })).toHaveValue(song.youtubeLink)

    expect(screen.getByRole('textbox', { name: /songsterr/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /songsterr/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /songsterr/i })).toHaveValue(song.songsterrLink)

    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
    await user.hover(screen.getByRole('button', { name: /save/i }))
    expect(await screen.findByText(/need to make a change/i)).toBeInTheDocument()
  })

  it('should send update request when the links have changed', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()

    let capturedRequest: UpdateSongRequest
    server.use(
      http.put('/songs', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const [{ rerender }] = reduxRender(
      withToastify(<EditSongLinksModal opened={true} onClose={onClose} song={song} />)
    )

    const youtubeTextBox = screen.getByRole('textbox', { name: /youtube/i })
    const songsterrTextBox = screen.getByRole('textbox', { name: /songsterr/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.clear(youtubeTextBox)
    await user.click(youtubeTextBox)
    await user.paste(newYoutubeLink)

    await user.clear(songsterrTextBox)
    await user.click(songsterrTextBox)
    await user.paste(newSongsterrLink)
    expect(saveButton).not.toHaveAttribute('data-disabled')
    await user.click(saveButton)

    expect(capturedRequest).toStrictEqual({
      ...song,
      youtubeLink: newYoutubeLink,
      songsterrLink: newSongsterrLink
    })
    expect(onClose).toHaveBeenCalledOnce()

    expect(await screen.findByText(/song links updated/i)).toBeInTheDocument()

    rerender(
      <EditSongLinksModal
        opened={true}
        onClose={onClose}
        song={{ ...song, youtubeLink: newYoutubeLink, songsterrLink: newSongsterrLink }}
      />
    )
    expect(saveButton).toHaveAttribute('data-disabled', 'true')
  })

  it('should send update request when the links have changed to null', async () => {
    const user = userEvent.setup()

    let capturedRequest: UpdateSongRequest
    server.use(
      http.put('/songs', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(<EditSongLinksModal opened={true} onClose={() => {}} song={song} />)

    await user.clear(screen.getByRole('textbox', { name: /youtube/i }))
    await user.clear(screen.getByRole('textbox', { name: /songsterr/i }))
    await user.click(screen.getByRole('button', { name: /save/i }))

    expect(capturedRequest).toStrictEqual({
      ...song,
      youtubeLink: null,
      songsterrLink: null
    })
  })

  it('should keep the save button disabled when the links have not changed', async () => {
    const user = userEvent.setup()

    reduxRender(<EditSongLinksModal opened={true} onClose={() => {}} song={song} />)

    const youtubeTextBox = screen.getByRole('textbox', { name: /youtube/i })
    const songsterrTextBox = screen.getByRole('textbox', { name: /songsterr/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.clear(youtubeTextBox)
    await user.click(youtubeTextBox)
    await user.paste(newYoutubeLink)

    expect(saveButton).not.toHaveAttribute('data-disabled')
    await user.hover(saveButton)
    expect(screen.queryByText(/need to make a change/i)).not.toBeInTheDocument()

    await user.clear(songsterrTextBox)
    await user.click(songsterrTextBox)
    await user.paste(newSongsterrLink)
    expect(saveButton).not.toHaveAttribute('data-disabled')

    await user.clear(youtubeTextBox)
    await user.click(youtubeTextBox)
    await user.paste(song.youtubeLink)
    expect(saveButton).not.toHaveAttribute('data-disabled')

    await user.clear(songsterrTextBox)
    await user.click(songsterrTextBox)
    await user.paste(song.songsterrLink)

    expect(saveButton).toHaveAttribute('data-disabled', 'true')
  })

  it('should validate links fields', async () => {
    const user = userEvent.setup()

    reduxRender(<EditSongLinksModal opened={true} onClose={() => {}} song={song} />)

    const youtubeTextBox = screen.getByRole('textbox', { name: /youtube/i })
    const songsterrTextBox = screen.getByRole('textbox', { name: /songsterr/i })

    // invalidate Youtube link
    await user.clear(youtubeTextBox)
    await user.type(youtubeTextBox, 'some invalid link')
    act(() => youtubeTextBox.blur())
    expect(youtubeTextBox).toBeInvalid()

    await user.clear(youtubeTextBox)
    expect(youtubeTextBox).not.toBeInvalid()

    // invalidate Songsterr link
    await user.clear(songsterrTextBox)
    await user.type(songsterrTextBox, 'some songsterr invalid link')
    act(() => songsterrTextBox.blur())
    expect(songsterrTextBox).toBeInvalid()

    await user.clear(songsterrTextBox)
    expect(songsterrTextBox).not.toBeInvalid()
  })
})
