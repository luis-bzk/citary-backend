import { RoleDB } from '../../data/interfaces';
import { GetAllFiltersDto } from '../../domain/schemas/general';
import { CreateRoleDto, UpdateRoleDto } from '../../domain/schemas/role';

export abstract class RoleDataSource {
  abstract findRoleByName(name: string): Promise<RoleDB>;

  abstract createRole(createRoleDto: CreateRoleDto): Promise<RoleDB>;

  abstract findRoleById(id: number): Promise<RoleDB>;

  abstract findRoleByNameId(id: number, name: string): Promise<RoleDB>;

  abstract updateRole(updateRoleDto: UpdateRoleDto): Promise<RoleDB>;

  abstract getAllRoles(getAllRolesDto: GetAllFiltersDto): Promise<RoleDB[]>;

  abstract deleteRole(id: number): Promise<RoleDB>;

  abstract getRolesById(ids: number[]): Promise<RoleDB[]>;
}
