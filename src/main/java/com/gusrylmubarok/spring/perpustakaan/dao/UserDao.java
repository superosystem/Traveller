package com.gusrylmubarok.spring.perpustakaan.dao;

import com.gusrylmubarok.spring.perpustakaan.entity.User;
import org.springframework.stereotype.Repository;

import javax.persistence.criteria.CriteriaBuilder;
import javax.persistence.criteria.Predicate;
import javax.persistence.criteria.Root;
import java.util.List;


@Repository
public class UserDao extends BaseDao<User> {
    @Override
    public List<Predicate> predicates(User param, CriteriaBuilder builder, Root<User> root, boolean isCount) {
        List<Predicate> predicates = super.predicates(param, builder, root, isCount);
        if (param != null) {
            if (param.getUsername() != null) {
                predicates.add(builder.equal(root.get("username"), param.getUsername()));
            }
        }
        return predicates;
    }
}
