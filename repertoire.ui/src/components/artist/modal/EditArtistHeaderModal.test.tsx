import { emptyArtist, reduxRender, withToastify } from '../../../test-utils.tsx'
import EditArtistHeaderModal from './EditArtistHeaderModal.tsx'
import Artist from '../../../types/models/Artist.ts'
import { act, screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { UpdateArtistRequest } from '../../../types/requests/ArtistRequests.ts'

describe('Edit Artist Header Modal', () => {
  const artist: Artist = {
    ...emptyArtist,
    id: '1',
    name: 'Artist 1',
    imageUrl: 'some-image.png'
  }

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<EditArtistHeaderModal opened={true} onClose={() => {}} artist={artist} />)

    expect(screen.getByRole('dialog', { name: /edit artist header/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /edit artist header/i })).toBeInTheDocument()

    expect(screen.getByRole('img', { name: 'image-preview' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: 'image-preview' })).toHaveAttribute(
      'src',
      artist.imageUrl
    )

    expect(screen.getByRole('textbox', { name: /name/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue(artist.name)

    expect(screen.getByRole('checkbox', { name: 'is-band' })).toBeInTheDocument()
    if (artist.isBand) {
      expect(screen.getByRole('checkbox', { name: 'is-band' })).toBeChecked()
    } else {
      expect(screen.getByRole('checkbox', { name: 'is-band' })).not.toBeChecked()
    }

    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
    await user.hover(screen.getByRole('button', { name: /save/i }))
    expect(await screen.findByText(/need to make a change/i)).toBeInTheDocument()
  })

  it('should send only edit request when the image is unchanged', async () => {
    const user = userEvent.setup()

    const newName = 'New Artist'
    const onClose = vitest.fn()

    let capturedRequest: UpdateArtistRequest
    server.use(
      http.put('/artists', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateArtistRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditArtistHeaderModal opened={true} onClose={onClose} artist={artist} />)
    )

    const nameField = screen.getByRole('textbox', { name: /name/i })
    const isBandCheckbox = screen.getByRole('checkbox', { name: 'is-band' })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.clear(nameField)
    await user.type(nameField, newName)
    await user.click(isBandCheckbox)
    await user.click(saveButton)

    expect(await screen.findByText(/artist updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      id: artist.id,
      name: newName,
      isBand: true
    })
  })

  it('should send edit request and save image request when the image is replaced', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const onClose = vitest.fn()

    let capturedRequest: UpdateArtistRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.put('/artists', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateArtistRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.put('/artists/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditArtistHeaderModal opened={true} onClose={onClose} artist={artist} />)
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.upload(screen.getByTestId('upload-image-input'), newImage)
    await user.click(saveButton)

    expect(await screen.findByText(/artist updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      id: artist.id,
      isBand: artist.isBand,
      name: artist.name
    })
    expect(capturedSaveImageFormData.get('id')).toBe(artist.id)
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
  })

  it('should send edit request and save image request when the image is first added', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const onClose = vitest.fn()

    let capturedRequest: UpdateArtistRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.put('/artists', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateArtistRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.put('/artists/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(
        <EditArtistHeaderModal
          opened={true}
          onClose={onClose}
          artist={{ ...artist, imageUrl: undefined }}
        />
      )
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)
    await user.click(saveButton)

    expect(await screen.findByText(/artist updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      id: artist.id,
      isBand: artist.isBand,
      name: artist.name
    })
    expect(capturedSaveImageFormData.get('id')).toBe(artist.id)
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
  })

  it('should send edit request and delete image request', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()

    let capturedRequest: UpdateArtistRequest
    server.use(
      http.put('/artists', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateArtistRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.delete(`/artists/images/${artist.id}`, () => {
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditArtistHeaderModal opened={true} onClose={onClose} artist={artist} />)
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.click(screen.getByRole('button', { name: 'remove-image' }))
    await user.click(saveButton)

    expect(await screen.findByText(/artist updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      id: artist.id,
      isBand: artist.isBand,
      name: artist.name
    })
  })

  it('should disable the save button when no changes are made', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })

    reduxRender(<EditArtistHeaderModal opened={true} onClose={() => {}} artist={artist} />)

    const nameField = screen.getByRole('textbox', { name: /name/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    // change image
    await user.upload(screen.getByTestId('upload-image-input'), newImage)
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset image
    await user.click(screen.getByRole('button', { name: 'reset-image' }))
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change name
    await user.type(nameField, '1')
    act(() => nameField.blur())
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset name
    await user.clear(nameField)
    await user.type(nameField, artist.name)
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // remove image
    await user.click(screen.getByRole('button', { name: 'remove-image' }))
    expect(saveButton).not.toHaveAttribute('data-disabled')
  })

  it('should validate the name textbox', async () => {
    const user = userEvent.setup()

    reduxRender(<EditArtistHeaderModal opened={true} onClose={() => {}} artist={artist} />)

    const nameField = screen.getByRole('textbox', { name: /name/i })
    await user.clear(nameField)
    expect(nameField).toBeInvalid()
  })
})
