package com.gusrylmubarok.spring.perpustakaan.service;

import com.gusrylmubarok.spring.perpustakaan.dao.BaseDao;
import com.gusrylmubarok.spring.perpustakaan.entity.BaseEntity;
import org.springframework.transaction.annotation.Transactional;

import java.util.Collection;

public abstract class BaseService<T extends BaseEntity<T>> {
    protected abstract BaseDao<T> getDAO();

    @Transactional(readOnly = true)
    public T findOne(T param) {
        return getDAO().findOne(param);
    }

    @Transactional(readOnly = true)
    public Collection<T> find(T param, int offset, int limit) {
        return getDAO().find(param, offset, limit);
    }

    @Transactional(readOnly = true)
    public Long count(T param) {
        return getDAO().count(param);
    }

    @Transactional
    public T save(T entity) {
        return getDAO().save(entity);
    }

    @Transactional
    public T update(T entity) {
        if (entity.getId() != null) {
            return getDAO().update(entity);
        }
        return null;
    }

    @Transactional
    public boolean delete(Long id) {
        return getDAO().delete(getDAO().findReference(id)) != null;
    }
}

