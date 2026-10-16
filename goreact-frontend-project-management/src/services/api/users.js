import network from '@/utils/network';

const users = {
  async findAll(params) {
    return network.get(`/api/v1/users/page`, {
      params,
    });
  },
};

export default users;
