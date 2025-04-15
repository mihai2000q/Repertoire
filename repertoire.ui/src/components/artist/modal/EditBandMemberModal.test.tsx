import { reduxRender, withToastify } from '../../../test-utils.tsx'
import EditBandMemberModal from './EditBandMemberModal.tsx'
import { act, screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { BandMember, BandMemberRole } from '../../../types/models/Artist.ts'
import { UpdateBandMemberRequest } from '../../../types/requests/ArtistRequests.ts'

describe('Edit Band Member Header Modal', () => {
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

  const bandMember: BandMember = {
    id: '1',
    name: 'Member 1',
    imageUrl: 'some-image.png',
    color: '#456',
    roles: [bandMemberRoles[0]]
  }

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<EditBandMemberModal opened={true} onClose={() => {}} bandMember={bandMember} />)

    expect(screen.getByRole('dialog', { name: /edit band member/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /edit band member/i })).toBeInTheDocument()

    expect(screen.getByRole('img', { name: 'image-preview' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: 'image-preview' })).toHaveAttribute(
      'src',
      bandMember.imageUrl
    )

    expect(screen.getByRole('textbox', { name: /name/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue(bandMember.name)

    expect(screen.getByRole('button', { name: 'color-input' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'color-input' })).toHaveStyle(
      `backgroundColor: '${bandMember.color}'`
    )

    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
    await user.hover(screen.getByRole('button', { name: /save/i }))
    expect(await screen.findByText(/need to make a change/i)).toBeInTheDocument()
  })

  it('should send only edit request when the image is unchanged', async () => {
    const user = userEvent.setup()

    const newName = 'New Member'
    const newRoles = [bandMemberRoles[1]]
    const onClose = vitest.fn()

    let capturedRequest: UpdateBandMemberRequest
    server.use(
      http.put('/artists/band-members', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateBandMemberRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditBandMemberModal opened={true} onClose={onClose} bandMember={bandMember} />)
    )

    const nameField = screen.getByRole('textbox', { name: /name/i })
    const rolesField = screen.getByRole('textbox', { name: /roles/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.clear(nameField)
    await user.type(nameField, newName)
    await user.click(rolesField)
    for (const role of [...bandMember.roles, ...newRoles]) {
      await user.click(screen.getByRole('option', { name: role.name })) // remove old roles and add new ones
    }
    act(() => rolesField.blur())

    await user.click(saveButton)

    expect(await screen.findByText(`${newName} updated!`)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      id: bandMember.id,
      name: newName,
      color: bandMember.color,
      roleIds: newRoles.map((role) => role.id)
    })
  })

  it('should send only save image request when the image is replaced', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const onClose = vitest.fn()

    let capturedSaveImageFormData: FormData
    server.use(
      http.put('/artists/band-members/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditBandMemberModal opened={true} onClose={onClose} bandMember={bandMember} />)
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.click(screen.getByRole('button', { name: 'image-options' }))
    await user.upload(screen.getByTestId('upload-image-input'), newImage)
    await user.click(saveButton)

    expect(await screen.findByText(`${bandMember.name} updated!`)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedSaveImageFormData.get('id')).toBe(bandMember.id)
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
  })

  it('should send only save image request when the image is first added', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const onClose = vitest.fn()

    let capturedSaveImageFormData: FormData
    server.use(
      http.put('/artists/band-members/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(
        <EditBandMemberModal
          opened={true}
          onClose={onClose}
          bandMember={{ ...bandMember, imageUrl: undefined }}
        />
      )
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)
    await user.click(saveButton)

    expect(await screen.findByText(`${bandMember.name} updated!`)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedSaveImageFormData.get('id')).toBe(bandMember.id)
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
  })

  it('should send only delete image request when the image is removed', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()

    server.use(
      http.delete(`/artists/band-members/images/${bandMember.id}`, () => {
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditBandMemberModal opened={true} onClose={onClose} bandMember={bandMember} />)
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.click(screen.getByRole('button', { name: 'image-options' }))
    await user.click(screen.getByRole('menuitem', { name: /remove image/i }))
    await user.click(saveButton)

    expect(await screen.findByText(`${bandMember.name} updated!`)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should send edit request and save image request when both have changed', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })
    const newName = 'New Member'
    const onClose = vitest.fn()

    let capturedRequest: UpdateBandMemberRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.put('/artists/band-members', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateBandMemberRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.put('/artists/band-members/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(<EditBandMemberModal opened={true} onClose={onClose} bandMember={bandMember} />)
    )

    const nameField = screen.getByRole('textbox', { name: /name/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await user.clear(nameField)
    await user.type(nameField, newName)
    await user.click(screen.getByRole('button', { name: 'image-options' }))
    await user.upload(screen.getByTestId('upload-image-input'), newImage)
    await user.click(saveButton)

    expect(await screen.findByText(`${newName} updated!`)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      id: bandMember.id,
      name: newName,
      color: bandMember.color,
      roleIds: bandMember.roles.map(r => r.id)
    })
    expect(capturedSaveImageFormData.get('id')).toBe(bandMember.id)
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
  })

  it('should disable the save button when no changes are made', async () => {
    const user = userEvent.setup()

    const newImage = new File(['something'], 'image.png', { type: 'image/png' })

    reduxRender(<EditBandMemberModal opened={true} onClose={() => {}} bandMember={bandMember} />)

    const nameField = screen.getByRole('textbox', { name: /name/i })
    const rolesField = screen.getByRole('textbox', { name: /roles/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    // change image
    await user.click(screen.getByRole('button', { name: 'image-options' }))
    await user.upload(screen.getByTestId('upload-image-input'), newImage)
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset image
    await user.click(screen.getByRole('button', { name: 'image-options' }))
    await user.click(screen.getByRole('menuitem', { name: /reset image/i }))
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change name
    await user.type(nameField, '1')
    act(() => nameField.blur())
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset name
    await user.clear(nameField)
    await user.type(nameField, bandMember.name)
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change roles
    await user.click(rolesField)
    for (const role of bandMemberRoles) {
      await user.click(screen.getByRole('option', { name: role.name }))
    }
    act(() => rolesField.blur())
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset roles
    await user.click(rolesField)
    for (const role of bandMemberRoles) {
      await user.click(screen.getByRole('option', { name: role.name }))
    }
    act(() => rolesField.blur())
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // remove image
    await user.click(screen.getByRole('button', { name: 'image-options' }))
    await user.click(screen.getByRole('menuitem', { name: /remove image/i }))
    expect(saveButton).not.toHaveAttribute('data-disabled')
  })

  it('should validate the name and roles textbox', async () => {
    const user = userEvent.setup()

    reduxRender(<EditBandMemberModal opened={true} onClose={() => {}} bandMember={bandMember} />)

    const nameField = screen.getByRole('textbox', { name: /name/i })
    const rolesField = screen.getByRole('textbox', { name: /roles/i })

    await user.clear(nameField)
    expect(nameField).toBeInvalid()

    await user.type(nameField, '1')

    await user.click(rolesField)
    for (const role of bandMember.roles) {
      await user.click(screen.getByRole('option', { name: role.name }))
    }
    expect(rolesField).toBeInvalid()
  })
})
