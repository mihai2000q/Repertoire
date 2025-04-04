import { emptySong, reduxRender, withToastify } from '../../../test-utils.tsx'
import Song from '../../../types/models/Song.ts'
import { setupServer } from 'msw/node'
import { screen } from '@testing-library/react'
import EditSongDescriptionModal from './EditSongDescriptionModal.tsx'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { UpdateSongRequest } from '../../../types/requests/SongRequests.ts'

describe('Edit Song Description Modal', () => {
  const song: Song = {
    ...emptySong,
    id: 'some-id',
    description: 'This is a description'
  }

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<EditSongDescriptionModal opened={true} onClose={() => {}} song={song} />)

    expect(screen.getByRole('dialog', { name: /edit song description/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /edit song description/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /description/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /description/i })).toHaveValue(song.description)
    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
    await user.hover(screen.getByRole('button', { name: /save/i }))
    expect(await screen.findByText(/need to make a change/i)).toBeInTheDocument()
  })

  it('should send update request when the description has changed', async () => {
    const user = userEvent.setup()

    const newDescription = 'This is a new description of the song'
    const onClose = vitest.fn()

    let capturedRequest: UpdateSongRequest
    server.use(
      http.put('/songs', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const [{ rerender }] = reduxRender(
      withToastify(<EditSongDescriptionModal opened={true} onClose={onClose} song={song} />)
    )

    await user.clear(screen.getByRole('textbox', { name: /description/i }))
    await user.type(screen.getByRole('textbox', { name: /description/i }), newDescription)
    expect(screen.getByRole('button', { name: /save/i })).not.toHaveAttribute('data-disabled')
    await user.click(screen.getByRole('button', { name: /save/i }))

    expect(capturedRequest).toStrictEqual({
      ...song,
      description: newDescription
    })
    expect(onClose).toHaveBeenCalledOnce()

    expect(await screen.findByText(/song description updated/i)).toBeInTheDocument()

    rerender(
      <EditSongDescriptionModal
        opened={true}
        onClose={onClose}
        song={{ ...song, description: newDescription }}
      />
    )
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
  })

  it('should send update request when the description has changed to empty', async () => {
    const user = userEvent.setup()

    let capturedRequest: UpdateSongRequest
    server.use(
      http.put('/songs', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditSongDescriptionModal opened={true} onClose={() => {}} song={song} />)
    )

    await user.clear(screen.getByRole('textbox', { name: /description/i }))
    await user.click(screen.getByRole('button', { name: /save/i }))

    expect(capturedRequest).toStrictEqual({
      ...song,
      description: ''
    })
  })

  it('should keep the save button disabled when the description has not changed', async () => {
    const user = userEvent.setup()

    reduxRender(<EditSongDescriptionModal opened={true} onClose={() => {}} song={song} />)

    await user.clear(screen.getByRole('textbox', { name: /description/i }))
    await user.type(screen.getByRole('textbox', { name: /description/i }), song.description + '1')
    expect(screen.getByRole('button', { name: /save/i })).not.toHaveAttribute('data-disabled')

    await user.clear(screen.getByRole('textbox', { name: /description/i }))
    await user.type(screen.getByRole('textbox', { name: /description/i }), song.description)
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
  })
})
