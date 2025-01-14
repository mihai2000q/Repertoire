import { reduxRender, withToastify } from '../../../test-utils.tsx'
import EditSongHeaderModal from './EditSongHeaderModal.tsx'
import Song from '../../../types/models/Song.ts'
import { act, screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import dayjs from 'dayjs'
import { http, HttpResponse } from 'msw'
import { UpdateSongRequest } from '../../../types/requests/SongRequests.ts'

describe('Edit Song Header Modal', () => {
  const song: Song = {
    id: '1',
    title: 'Song 1',
    releaseDate: '2024-12-11T23:00:00.000Z',
    imageUrl: 'some-image.png',
    confidence: 0,
    isRecorded: false,
    progress: 0,
    rehearsals: 0,
    sections: [],
    description: '',
    createdAt: '',
    updatedAt: ''
  }

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<EditSongHeaderModal opened={true} onClose={() => {}} song={song} />)

    expect(screen.getByRole('dialog', { name: /edit song header/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /edit song header/i })).toBeInTheDocument()

    expect(screen.getByRole('img', { name: 'image-preview' })).toBeInTheDocument()

    expect(screen.getByRole('textbox', { name: /title/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /title/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /title/i })).toHaveValue(song.title)

    expect(screen.getByRole('button', { name: /release date/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /release date/i })).toHaveTextContent(
      dayjs(song.releaseDate).format('D MMMM YYYY')
    )

    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
    await user.hover(screen.getByRole('button', { name: /save/i }))
    expect(await screen.findByText(/need to make a change/i)).toBeInTheDocument()
  })

  it('should send only edit request when the image is unchanged', async () => {
    const user = userEvent.setup()

    const newTitle = 'New Song'
    const newReleaseDate = dayjs('2024-12-24')
    const onClose = vitest.fn()

    let capturedRequest: UpdateSongRequest
    server.use(
      http.put('/songs', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(withToastify(<EditSongHeaderModal opened={true} onClose={onClose} song={song} />))

    const titleField = screen.getByRole('textbox', { name: /title/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.clear(titleField)
    await user.type(titleField, newTitle)
    await user.click(screen.getByRole('button', { name: /release date/i }))
    await user.click(screen.getByText(newReleaseDate.date()))
    await user.click(saveButton)

    expect(await screen.findByText(/song header updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      ...song,
      id: song.id,
      title: newTitle,
      releaseDate: newReleaseDate.toISOString()
    })
  })

  it('should send edit request and save image request when the image is replaced', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const onClose = vitest.fn()

    let capturedRequest: UpdateSongRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.put('/songs', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.put('/songs/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(withToastify(<EditSongHeaderModal opened={true} onClose={onClose} song={song} />))

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.upload(screen.getByTestId('upload-image-input'), newImage)
    await user.click(saveButton)

    expect(await screen.findByText(/song header updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      ...song,
      id: song.id,
      title: song.title,
      releaseDate: dayjs(song.releaseDate).toISOString()
    })
    expect(capturedSaveImageFormData.get('id')).toBe(song.id)
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
  })

  it('should send edit request and save image request when the image is first added', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const onClose = vitest.fn()

    let capturedRequest: UpdateSongRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.put('/songs', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.put('/songs/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(
        <EditSongHeaderModal
          opened={true}
          onClose={onClose}
          song={{ ...song, imageUrl: undefined }}
        />
      )
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)
    await user.click(saveButton)

    expect(await screen.findByText(/song header updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toEqual({
      ...{ ...song, imageUrl: undefined },
      id: song.id,
      title: song.title,
      releaseDate: dayjs(song.releaseDate).toISOString()
    })
    expect(capturedSaveImageFormData.get('id')).toBe(song.id)
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
  })

  it('should send edit request and delete image request when the image is first added', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()

    let capturedRequest: UpdateSongRequest
    server.use(
      http.put('/songs', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.delete(`/songs/images/${song.id}`, () => {
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(withToastify(<EditSongHeaderModal opened={true} onClose={onClose} song={song} />))

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.click(screen.getByRole('button', { name: 'remove-image' }))
    await user.click(saveButton)

    expect(await screen.findByText(/song header updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      ...song,
      id: song.id,
      title: song.title,
      releaseDate: dayjs(song.releaseDate).toISOString()
    })
  })

  it('should disable the save button when no changes are made', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const day = new Date(song.releaseDate).getDate()

    reduxRender(<EditSongHeaderModal opened={true} onClose={() => {}} song={song} />)

    const titleField = screen.getByRole('textbox', { name: /title/i })
    const releaseDateField = screen.getByRole('button', { name: /release date/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    // change image
    await user.upload(screen.getByTestId('upload-image-input'), newImage)
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset image
    await user.click(screen.getByRole('button', { name: 'reset-image' }))
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change title
    await user.type(titleField, '1')
    act(() => titleField.blur())
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset title
    await user.clear(titleField)
    await user.type(titleField, song.title)
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change release date
    await user.click(releaseDateField)
    await user.click(screen.getByText((day + 1).toString()))
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset release date
    await user.click(releaseDateField)
    await user.click(screen.getByText(day.toString()))
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // remove image
    await user.click(screen.getByRole('button', { name: 'remove-image' }))
    expect(saveButton).not.toHaveAttribute('data-disabled')
  })

  it('should validate the title textbox', async () => {
    const user = userEvent.setup()

    reduxRender(<EditSongHeaderModal opened={true} onClose={() => {}} song={song} />)

    const titleField = screen.getByRole('textbox', { name: /title/i })
    await user.clear(titleField)
    expect(titleField).toBeInvalid()
  })
})
