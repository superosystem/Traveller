package com.gusrylmubarok.spring.perpustakaan.controller;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.gusrylmubarok.spring.perpustakaan.common.RestResult;
import com.gusrylmubarok.spring.perpustakaan.common.StatusCode;
import com.gusrylmubarok.spring.perpustakaan.entity.Book;
import com.gusrylmubarok.spring.perpustakaan.service.BookService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.util.ArrayList;

@RestController
@RequestMapping("/books")
@PreAuthorize("permitAll()")
public class BookController extends BaseController {
    @Autowired
    private BookService service;

    @GetMapping
    public RestResult get(@RequestParam(value = "param", required = false) String param,
                          @RequestParam(value = "offset") int offset,
                          @RequestParam(value = "limit") int limit) throws JsonProcessingException {
        Book book = param != null ? new ObjectMapper().readValue(param, Book.class) : null;
        long rows = service.count(book);
        return new RestResult(rows > 0 ? service.find(book, offset, limit) : new ArrayList<>(), rows);
    }

    @PostMapping
    public RestResult save(@RequestBody Book param) {
        param = service.save(param);
        return new RestResult(param, param != null ? StatusCode.SAVE_SUCCESS : StatusCode.SAVE_FAILED);
    }

    @PutMapping
    public RestResult update(@RequestBody Book book) {
        book = service.update(book);
        return new RestResult(book, book != null ? StatusCode.UPDATE_SUCCESS : StatusCode.UPDATE_FAILED);
    }

    @DeleteMapping(value = "{id}")
    public RestResult delete(@PathVariable Long id) {
        return new RestResult(service.delete(id) ? StatusCode.DELETE_SUCCESS : StatusCode.DELETE_FAILED);
    }
}
