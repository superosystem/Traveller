package com.gusrylmubarok.spring.perpustakaan.service;

import com.gusrylmubarok.spring.perpustakaan.dao.BaseDao;
import com.gusrylmubarok.spring.perpustakaan.dao.BookDao;
import com.gusrylmubarok.spring.perpustakaan.entity.Book;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class BookService extends BaseService<Book> {
    @Autowired
    private BookDao dao;
    @Override
    protected BaseDao<Book> getDAO() {
        return dao;
    }
}
