import { emptyPlaylist, reduxRender, withToastify } from '../../../test-utils.tsx'
import EditPlaylistHeaderModal from './EditPlaylistHeaderModal.tsx'
import Playlist from '../../../types/models/Playlist.ts'
import { act, screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { UpdatePlaylistRequest } from '../../../types/requests/PlaylistRequests.ts'

describe('Edit Playlist Header Modal', () => {
  const playlist: Playlist = {
    ...emptyPlaylist,
    id: '1',
    title: 'Playlist 1',
    description: 'This is a description',
    imageUrl: 'some-image.png'
  }

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<EditPlaylistHeaderModal opened={true} onClose={() => {}} playlist={playlist} />)

    expect(screen.getByRole('dialog', { name: /edit playlist header/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /edit playlist header/i })).toBeInTheDocument()

    expect(screen.getByRole('img', { name: 'image-preview' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: 'image-preview' })).toHaveAttribute(
      'src',
      playlist.imageUrl
    )

    expect(screen.getByRole('textbox', { name: /title/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /title/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /title/i })).toHaveValue(playlist.title)

    expect(screen.getByRole('textbox', { name: /description/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /description/i })).toHaveValue(playlist.description)

    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
    await user.hover(screen.getByRole('button', { name: /save/i }))
    expect(await screen.findByText(/need to make a change/i)).toBeInTheDocument()
  })

  it('should send only edit request when the image is unchanged', async () => {
    const user = userEvent.setup()

    const newTitle = 'New Playlist'
    const newDescription = 'This is the new description'
    const onClose = vitest.fn()

    let capturedRequest: UpdatePlaylistRequest
    server.use(
      http.put('/playlists', async (req) => {
        capturedRequest = (await req.request.json()) as UpdatePlaylistRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditPlaylistHeaderModal opened={true} onClose={onClose} playlist={playlist} />)
    )

    const titleField = screen.getByRole('textbox', { name: /title/i })
    const descriptionField = screen.getByRole('textbox', { name: /description/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.clear(titleField)
    await user.type(titleField, newTitle)

    await user.clear(descriptionField)
    await user.type(descriptionField, newDescription)

    await user.click(saveButton)

    expect(await screen.findByText(/playlist updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      id: playlist.id,
      title: newTitle,
      description: newDescription
    })
  })

  it('should send only save image request when the image is replaced', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const onClose = vitest.fn()

    let capturedSaveImageFormData: FormData
    server.use(
      http.put('/playlists/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditPlaylistHeaderModal opened={true} onClose={onClose} playlist={playlist} />)
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.upload(screen.getByTestId('upload-image-input'), newImage)
    await user.click(saveButton)

    expect(await screen.findByText(/playlist updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedSaveImageFormData.get('id')).toBe(playlist.id)
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
  })

  it('should send only save image request when the image is first added', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const onClose = vitest.fn()

    let capturedSaveImageFormData: FormData
    server.use(
      http.put('/playlists/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(
        <EditPlaylistHeaderModal
          opened={true}
          onClose={onClose}
          playlist={{ ...playlist, imageUrl: undefined }}
        />
      )
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)
    await user.click(saveButton)

    expect(await screen.findByText(/playlist updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedSaveImageFormData.get('id')).toBe(playlist.id)
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
  })

  it('should send only delete image request when the image is removed', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()

    server.use(
      http.delete(`/playlists/images/${playlist.id}`, () => {
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditPlaylistHeaderModal opened={true} onClose={onClose} playlist={playlist} />)
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.click(screen.getByRole('button', { name: 'remove-image' }))
    await user.click(saveButton)

    expect(await screen.findByText(/playlist updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should send edit request and save image request when both have changed', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const newTitle = 'New Playlist'
    const onClose = vitest.fn()

    let capturedRequest: UpdatePlaylistRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.put('/playlists', async (req) => {
        capturedRequest = (await req.request.json()) as UpdatePlaylistRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.put('/playlists/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditPlaylistHeaderModal opened={true} onClose={onClose} playlist={playlist} />)
    )

    const titleField = screen.getByRole('textbox', { name: /title/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.clear(titleField)
    await user.type(titleField, newTitle)
    await user.upload(screen.getByTestId('upload-image-input'), newImage)
    await user.click(saveButton)

    expect(await screen.findByText(/playlist updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      id: playlist.id,
      title: newTitle,
      description: playlist.description
    })
    expect(capturedSaveImageFormData.get('id')).toBe(playlist.id)
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
  })

  it('should disable the save button when no changes are made', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })

    reduxRender(<EditPlaylistHeaderModal opened={true} onClose={() => {}} playlist={playlist} />)

    const titleField = screen.getByRole('textbox', { name: /title/i })
    const descriptionField = screen.getByRole('textbox', { name: /description/i })
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
    await user.type(titleField, playlist.title)
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change description
    await user.type(descriptionField, '1')
    act(() => descriptionField.blur())
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset description
    await user.clear(descriptionField)
    await user.type(descriptionField, playlist.description)
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // remove image
    await user.click(screen.getByRole('button', { name: 'remove-image' }))
    expect(saveButton).not.toHaveAttribute('data-disabled')
  })

  it('should validate the title textbox', async () => {
    const user = userEvent.setup()

    reduxRender(<EditPlaylistHeaderModal opened={true} onClose={() => {}} playlist={playlist} />)

    const titleField = screen.getByRole('textbox', { name: /title/i })
    await user.clear(titleField)
    expect(titleField).toBeInvalid()
  })
})
