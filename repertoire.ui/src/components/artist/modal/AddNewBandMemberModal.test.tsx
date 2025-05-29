import { reduxRender, withToastify } from '../../../test-utils.tsx'
import { setupServer } from 'msw/node'
import AddNewBandMemberModal from './AddNewBandMemberModal.tsx'
import { act, screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { CreateBandMemberRequest } from '../../../types/requests/ArtistRequests.ts'
import { BandMemberRole } from '../../../types/models/Artist.ts'

describe('Add New Band Member Modal', () => {
  const bandMemberRoles: BandMemberRole[] = [
    {
      id: '1',
      name: 'Guitarist'
    },
    {
      id: '2',
      name: 'Voice'
    },
    {
      id: '3',
      name: 'Drummer'
    }
  ]

  const handlers = [
    http.get('/artists/band-members/roles', async () => {
      return HttpResponse.json(bandMemberRoles)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', () => {
    reduxRender(<AddNewBandMemberModal opened={true} onClose={() => {}} artistId={''} />)

    expect(screen.getByRole('dialog', { name: /add new band member/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /add new band member/i })).toBeInTheDocument()
    expect(screen.getByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'color-input' })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /roles/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /submit/i })).toBeInTheDocument()
  })

  it('should send only create request when no image is uploaded - minimal', async () => {
    const user = userEvent.setup()

    const artistId = 'some-artist-id'
    const newName = 'New Member'
    const role = bandMemberRoles[1]

    const onClose = vitest.fn()

    let capturedRequest: CreateBandMemberRequest
    server.use(
      http.post('/artists/band-members', async (req) => {
        capturedRequest = (await req.request.json()) as CreateBandMemberRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<AddNewBandMemberModal opened={true} onClose={onClose} artistId={artistId} />)
    )

    await user.type(screen.getByRole('textbox', { name: /name/i }), newName)
    await user.click(screen.getByRole('textbox', { name: /roles/i }))
    await user.click(await screen.findByRole('option', { name: role.name }))
    act(() => screen.getByRole('textbox', { name: /roles/i }).blur())
    await user.click(screen.getByRole('button', { name: /submit/i }))

    await waitFor(() =>
      expect(capturedRequest).toStrictEqual({
        name: newName,
        roleIds: [role.id],
        artistId: artistId
      })
    )
    expect(onClose).toHaveBeenCalledOnce()

    expect(screen.getByText(`${newName} added!`))
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue('')
  })

  it('should send create request and save image request when the image is uploaded', async () => {
    const user = userEvent.setup()

    const artistId = 'some-artist-id'
    const newName = 'New Member'
    const role = bandMemberRoles[1]
    const newImage = new File(['something'], 'image.png', { type: 'image/png' })

    const onClose = vitest.fn()

    const returnedId = 'the-song-id'

    let capturedCreateRequest: CreateBandMemberRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.post('/artists/band-members', async (req) => {
        capturedCreateRequest = (await req.request.json()) as CreateBandMemberRequest
        return HttpResponse.json({ id: returnedId })
      }),
      http.put('/artists/band-members/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<AddNewBandMemberModal opened={true} onClose={onClose} artistId={artistId} />)
    )

    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)
    await user.type(screen.getByRole('textbox', { name: /name/i }), newName)
    await user.click(screen.getByRole('textbox', { name: /roles/i }))
    await user.click(await screen.findByRole('option', { name: role.name }))
    act(() => screen.getByRole('textbox', { name: /roles/i }).blur())
    await user.click(screen.getByRole('button', { name: /submit/i }))

    await waitFor(() =>
      expect(capturedCreateRequest).toStrictEqual({
        name: newName,
        roleIds: [role.id],
        artistId: artistId
      })
    )
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
    expect(capturedSaveImageFormData.get('id')).toBe(returnedId)

    expect(onClose).toHaveBeenCalledOnce()

    expect(screen.getByText(`${newName} added!`))
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue('')
    expect(screen.getByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
  })

  it('should remove the image on close', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })

    const onClose = vitest.fn()

    reduxRender(<AddNewBandMemberModal opened={true} onClose={onClose} artistId={''} />)

    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)
    expect(screen.getByRole('img', { name: 'image-preview' })).toBeInTheDocument()

    await user.click(screen.getAllByRole('button').find((b) => b.className.includes('CloseButton')))

    expect(await screen.findByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should validate name and roles', async () => {
    const user = userEvent.setup()

    reduxRender(<AddNewBandMemberModal opened={true} onClose={() => {}} artistId={''} />)

    const name = screen.getByRole('textbox', { name: /name/i })
    const roles = screen.getByRole('textbox', { name: /roles/i })

    expect(name).not.toBeInvalid()
    expect(roles).not.toBeInvalid()

    await user.click(screen.getByRole('button', { name: /submit/i }))
    expect(name).toBeInvalid()

    await user.type(name, 'something')
    expect(name).not.toBeInvalid()

    await user.click(screen.getByRole('button', { name: /submit/i }))
    expect(roles).toBeInvalid()

    await user.click(roles)
    await user.click((await screen.findAllByRole('option'))[0])
    expect(roles).not.toBeInvalid()

    await user.clear(name)
    act(() => name.blur())
    expect(name).toBeInvalid()
  })
})
