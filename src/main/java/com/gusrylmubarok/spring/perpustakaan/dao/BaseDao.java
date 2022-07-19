package com.gusrylmubarok.spring.perpustakaan.dao;

import com.gusrylmubarok.spring.perpustakaan.entity.BaseEntity;
import org.hibernate.Hibernate;

import javax.persistence.EntityManager;
import javax.persistence.NoResultException;
import javax.persistence.PersistenceContext;
import javax.persistence.TypedQuery;
import javax.persistence.criteria.CriteriaBuilder;
import javax.persistence.criteria.CriteriaQuery;
import javax.persistence.criteria.Predicate;
import javax.persistence.criteria.Root;
import java.lang.reflect.ParameterizedType;
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;

@SuppressWarnings({"unchecked", "ConstantConditions"})
public abstract class BaseDao<T extends BaseEntity<T>> {
    @PersistenceContext
    protected EntityManager entityManager;

    private Class<T> type;

    {
        this.type = (Class<T>) ((ParameterizedType) getClass().getGenericSuperclass()).getActualTypeArguments()[0];
    }

    public T findOne(T param) {
        CriteriaBuilder builder = entityManager.getCriteriaBuilder();
        CriteriaQuery<T> query = builder.createQuery(type);
        Root<T> root = query.from(type);
        query.orderBy(builder.asc(root.get("id")));

        return singleResult(query, predicates(param, builder, root, false));
    }

    public Collection<T> find(T param, int offset, int limit) {
        CriteriaBuilder builder = entityManager.getCriteriaBuilder();
        CriteriaQuery<T> query = builder.createQuery(type);
        Root<T> root = query.from(type);
        query.orderBy(builder.asc(root.get("id")));

        return listResult(query, predicates(param, builder, root, false), offset, limit);
    }

    public Long count(T param) {
        CriteriaBuilder builder = entityManager.getCriteriaBuilder();
        CriteriaQuery<Long> query = builder.createQuery(Long.class);
        Root<T> root =  query.from(type);
        query.select(builder.count(root));

        return singleResult(query, predicates(param, builder, root, true));
    }

    public T save(T entity) {
        if (entity != null && entity.getId() == null) {
            entityManager.persist(entity);
            return entity;
        }
        return entity;
    }

    public T update(T entity) {
        if((entity != null) && (entity.getId() != null)) {
            T reference = findReference(entity.getId());
            entity.setCreatedTime(reference.getCreatedTime());
            if (reference != null) {
                entityManager.merge(entity);
                return entity;
            }
        }
        return entity;
    }

    public T delete(T entity) {
        if (entity != null) {
            T reference = findReference(entity.getId());
            if (reference != null) {
                entityManager.remove(entity);
                return entity;
            }
        }
        return entity;
    }

    public T findReference(Long id) {
        return (T) Hibernate.unproxy(entityManager.getReference(type, id));
    }

    public List<Predicate> predicates(T param, CriteriaBuilder builder, Root<T> root, boolean isCount) {
        return new ArrayList<>();
    }

    public <I> I singleResult(CriteriaQuery<I> query, List<Predicate> predicates) {
        try {
            query.where(predicates.toArray(new Predicate[0]));
            return entityManager.createQuery(query).getSingleResult();
        } catch (NoResultException ignore) {
        }
        return null;
    }

    public List<T> listResult(CriteriaQuery<T> query, List<Predicate> predicates, int offset, int limit) {
        query.where(predicates.toArray(new Predicate[0]));
        query.distinct(true);
        TypedQuery<T> typeQuery = entityManager.createQuery(query);
        if (limit != Integer.MAX_VALUE) {
            typeQuery.setMaxResults(limit);
        }
        return typeQuery.setFirstResult(offset).getResultList();
    }
}
