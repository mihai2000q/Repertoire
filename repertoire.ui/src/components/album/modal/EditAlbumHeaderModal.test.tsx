import { reduxRender, withToastify } from '../../../test-utils.tsx'
import EditAlbumHeaderModal from './EditAlbumHeaderModal.tsx'
import Album from '../../../types/models/Album.ts'
import { act, screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import dayjs from 'dayjs'
import { http, HttpResponse } from 'msw'
import { UpdateAlbumRequest } from '../../../types/requests/AlbumRequests.ts'

describe('Edit Album Header Modal', () => {
  const album: Album = {
    id: '1',
    title: 'Album 1',
    songs: [],
    createdAt: '',
    updatedAt: '',
    releaseDate: '2024-12-12T00:00:00',
    imageUrl: 'some-image.png'
  }

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<EditAlbumHeaderModal opened={true} onClose={() => {}} album={album} />)

    expect(screen.getByRole('dialog', { name: /edit album header/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /edit album header/i })).toBeInTheDocument()

    expect(screen.getByRole('img', { name: 'image-preview' })).toBeInTheDocument()

    expect(screen.getByRole('textbox', { name: /title/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /title/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /title/i })).toHaveValue(album.title)

    expect(screen.getByRole('button', { name: /release date/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /release date/i })).toHaveTextContent(
      dayjs(album.releaseDate).format('D MMMM YYYY')
    )

    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
    await user.hover(screen.getByRole('button', { name: /save/i }))
    expect(await screen.findByText(/need to make a change/i)).toBeInTheDocument()
  })

  it('should send only edit request when the image is unchanged', async () => {
    const user = userEvent.setup()

    const newTitle = 'New Album'
    const newDay = '24'
    const newReleaseDate = `2024-12-${newDay}`
    const onClose = vitest.fn()

    let capturedRequest: UpdateAlbumRequest
    server.use(
      http.put('/albums', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateAlbumRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditAlbumHeaderModal opened={true} onClose={onClose} album={album} />)
    )

    const titleField = screen.getByRole('textbox', { name: /title/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.clear(titleField)
    await user.type(titleField, newTitle)
    await user.click(screen.getByRole('button', { name: /release date/i }))
    await user.click(screen.getByText(newDay))
    await user.click(saveButton)

    expect(await screen.findByText(/album updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      id: album.id,
      title: newTitle,
      releaseDate: dayjs(newReleaseDate).toISOString()
    })
  })

  it('should send edit request and save image request when the image is replaced', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const onClose = vitest.fn()

    let capturedRequest: UpdateAlbumRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.put('/albums', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateAlbumRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.put('/albums/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditAlbumHeaderModal opened={true} onClose={onClose} album={album} />)
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.upload(screen.getByTestId('upload-image-input'), newImage)
    await user.click(saveButton)

    expect(await screen.findByText(/album updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      id: album.id,
      title: album.title,
      releaseDate: dayjs(album.releaseDate).toISOString()
    })
    expect(capturedSaveImageFormData.get('id')).toBe(album.id)
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
  })

  it('should send edit request and save image request when the image is first added', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const onClose = vitest.fn()

    let capturedRequest: UpdateAlbumRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.put('/albums', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateAlbumRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.put('/albums/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(
        <EditAlbumHeaderModal
          opened={true}
          onClose={onClose}
          album={{ ...album, imageUrl: undefined }}
        />
      )
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)
    await user.click(saveButton)

    expect(await screen.findByText(/album updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      id: album.id,
      title: album.title,
      releaseDate: dayjs(album.releaseDate).toISOString()
    })
    expect(capturedSaveImageFormData.get('id')).toBe(album.id)
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
  })

  it('should send edit request and delete image request when the image is first added', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()

    let capturedRequest: UpdateAlbumRequest
    server.use(
      http.put('/albums', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateAlbumRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.delete(`/albums/images/${album.id}`, () => {
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditAlbumHeaderModal opened={true} onClose={onClose} album={album} />)
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.click(screen.getByRole('button', { name: 'remove-image' }))
    await user.click(saveButton)

    expect(await screen.findByText(/album updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      id: album.id,
      title: album.title,
      releaseDate: dayjs(album.releaseDate).toISOString()
    })
  })

  it('should disable the save button when no changes are made', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const day = new Date(album.releaseDate).getDate()

    reduxRender(<EditAlbumHeaderModal opened={true} onClose={() => {}} album={album} />)

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
    await user.type(titleField, album.title)
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

    reduxRender(<EditAlbumHeaderModal opened={true} onClose={() => {}} album={album} />)

    const titleField = screen.getByRole('textbox', { name: /title/i })
    await user.clear(titleField)
    expect(titleField).toBeInvalid()
  })
})
